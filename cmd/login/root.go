package login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/dbbackup-io/cli/pkg/config"
	"github.com/dbbackup-io/cli/pkg/logger"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

// Team represents a team from the API
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with DBBackup",
	Long:  "Start the authentication flow by opening the dashboard login page and starting a callback server",
	Run:   runLogin,
}

func runLogin(cmd *cobra.Command, args []string) {
	// Load configuration from environment variables with defaults
	const callbackPort = "9097"
	loginBaseURL := getEnvWithDefault("LOGIN_BASE_URL", "https://auth.dbbackup.io")

	callbackURL := "http://localhost:" + callbackPort
	loginURL := loginBaseURL + "/login?cli-login=true&cli-login-redirect=" + callbackURL

	// If already authenticated, we'll overwrite the existing token
	if config.IsAuthenticated() {
		logger.Debug("Existing authentication found, will overwrite with new token")
	}

	// Create a new HTTP mux for this server
	mux := http.NewServeMux()
	authComplete := make(chan map[string]interface{})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Debugf("Received callback: %s", r.URL.String())

		// Ignore favicon requests and other non-auth requests
		if r.URL.Path != "/" {
			logger.Debugf("Ignoring request for %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Parse token and teams from query parameters
		token := r.URL.Query().Get("token")
		teamsParam := r.URL.Query().Get("teams")

		logger.Debugf("Parsed token: %s, teams param: %s", token, teamsParam)

		// If no parameters at all, this might be a browser pre-request
		if token == "" && teamsParam == "" && len(r.URL.Query()) == 0 {
			logger.Debug("Ignoring request with no parameters")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Waiting for authentication...")
			return
		}

		if token == "" {
			fmt.Fprintf(w, "Authentication failed: No token received.")
			logger.Error("Authentication failed: No token received")
			authComplete <- nil
			return
		}

		if teamsParam == "" {
			fmt.Fprintf(w, "Authentication failed: No teams received.")
			logger.Error("Authentication failed: No teams received")
			authComplete <- nil
			return
		}

		// URL decode the teams parameter
		decodedTeams, err := url.QueryUnescape(teamsParam)
		if err != nil {
			fmt.Fprintf(w, "Authentication failed: Invalid teams data.")
			logger.Errorf("Failed to decode teams parameter: %v", err)
			authComplete <- nil
			return
		}

		// Parse teams JSON
		var teams []Team
		if err := json.Unmarshal([]byte(decodedTeams), &teams); err != nil {
			fmt.Fprintf(w, "Authentication failed: Invalid teams format.")
			logger.Errorf("Failed to parse teams JSON: %v", err)
			authComplete <- nil
			return
		}

		logger.Debugf("Parsed %d teams", len(teams))

		// Handle the callback
		fmt.Fprintf(w, "Authentication successful! You can close this window.")
		logger.Debug("Authentication callback processed successfully")
		authComplete <- map[string]interface{}{
			"token": token,
			"teams": teams,
		}
	})

	// Start callback server with dedicated mux
	server := &http.Server{
		Addr:    ":" + callbackPort,
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		logger.Debugf("Starting callback server on %s", callbackURL)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server error: %v", err)
		}
	}()

	// Open browser
	logger.Infof("Opening browser to: %s", "https://"+loginURL)
	if err := open.Run("https://" + loginURL); err != nil {
		logger.Warnf("Failed to open browser: %v", err)
		logger.Infof("Please manually open: https://%s", loginURL)
	}

	// Wait for authentication or timeout
	var authData map[string]interface{}
	select {
	case authData = <-authComplete:
		logger.Debug("Authentication callback received")
	case <-time.After(5 * time.Minute):
		logger.Warn("Authentication timed out after 5 minutes")
	}

	// Save auth data if authentication was successful
	if authData != nil {
		token, ok := authData["token"].(string)
		if !ok {
			logger.Error("Invalid token in authentication data")
			// Shutdown server before returning
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			return
		}

		teamsInterface, ok := authData["teams"]
		if !ok {
			logger.Error("No teams in authentication data")
			// Shutdown server before returning
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			return
		}

		teams, ok := teamsInterface.([]Team)
		if !ok {
			logger.Error("Invalid teams format in authentication data")
			// Shutdown server before returning
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			return
		}

		// Let user select team immediately
		selectedTeam, err := selectTeam(teams)
		if err != nil {
			logger.Errorf("Failed to select team: %v", err)
			// Shutdown server before returning
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			return
		}

		teamId := strconv.Itoa(selectedTeam.ID)
		logger.Debugf("User selected team: %s (ID: %s)", selectedTeam.Name, teamId)

		if err := config.SetAuthData(token, teamId); err != nil {
			logger.Errorf("Failed to save authentication data: %v", err)
			// Shutdown server before returning
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			return
		}
		logger.Infof("Authentication successful! Connected to team: %s", selectedTeam.Name)
	}

	// Shutdown server at the end
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

// selectTeam prompts the user to select a team using huh
func selectTeam(teams []Team) (*Team, error) {
	if len(teams) == 0 {
		return nil, fmt.Errorf("no teams available")
	}

	// If only one team, select it automatically
	if len(teams) == 1 {
		logger.Infof("Only one team available: %s", teams[0].Name)
		return &teams[0], nil
	}

	// Create options for huh select
	var options []huh.Option[int]
	for _, team := range teams {
		options = append(options, huh.NewOption(team.Name, team.ID))
	}

	var selectedTeamID int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select a team to connect to:").
				Options(options...).
				Value(&selectedTeamID),
		),
	)

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("team selection failed: %w", err)
	}

	// Find the selected team
	for _, team := range teams {
		if team.ID == selectedTeamID {
			return &team, nil
		}
	}

	return nil, fmt.Errorf("selected team not found")
}

// getEnvWithDefault returns the value of an environment variable or a default value if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
