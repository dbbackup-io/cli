package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dbbackup-io/cli/pkg/config"
	"github.com/dbbackup-io/cli/pkg/logger"
)

// Client represents the API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
	TeamID     string
}

// BackupJob represents a backup job from the API
type BackupJob struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	Schedule             string    `json:"schedule"`
	ExecutionMode        string    `json:"execution_mode"`
	ServerID             *int      `json:"server_id"`
	DatabaseSourceID     int       `json:"database_source_id"`
	StorageDestinationID int       `json:"storage_destination_id"`
	RetentionDays        int       `json:"retention_days"`
	Compression          string    `json:"compression"`
	Status               string    `json:"status"`
	Tags                 *string   `json:"tags"`
	TeamID               int       `json:"team_id"`
	WorkspaceID          int       `json:"workspace_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Server represents a server from the API
type Server struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	SSHHost     *string    `json:"ssh_host"`
	SSHPort     int        `json:"ssh_port"`
	SSHUser     string     `json:"ssh_user"`
	Workdir     string     `json:"workdir"`
	Status      string     `json:"status"`
	LastSeenAt  *time.Time `json:"last_seen_at"`
	Tags        *string    `json:"tags"`
	TeamID      int        `json:"team_id"`
	WorkspaceID int        `json:"workspace_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DatabaseSource represents a database source from the API
type DatabaseSource struct {
	ID          int       `json:"id"`
	Engine      string    `json:"engine"`
	Host        string    `json:"host"`
	Port        int       `json:"port"`
	DBName      *string   `json:"db_name"`
	Username    *string   `json:"username"`
	Password    *string   `json:"password"`
	ConfigJSON  *string   `json:"config_json"`
	Tags        *string   `json:"tags"`
	Status      string    `json:"status"`
	TeamID      int       `json:"team_id"`
	WorkspaceID int       `json:"workspace_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StorageDestination represents a storage destination from the API
type StorageDestination struct {
	ID          int       `json:"id"`
	Driver      string    `json:"driver"`
	Region      *string   `json:"region"`
	Bucket      string    `json:"bucket"`
	PathPrefix  *string   `json:"path_prefix"`
	ConfigJSON  *string   `json:"config_json"`
	IsDefault   bool      `json:"is_default"`
	Tags        *string   `json:"tags"`
	TeamID      int       `json:"team_id"`
	WorkspaceID int       `json:"workspace_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BackupRun represents a backup run from the API
type BackupRun struct {
	ID          int        `json:"id"`
	StartedAt   time.Time  `json:"started_at"`
	BackupJobID int        `json:"backup_job_id"`
	FinishedAt  *time.Time `json:"finished_at"`
	Status      string     `json:"status"`
	RetryCount  int        `json:"retry_count"`
	SizeBytes   *int64     `json:"size_bytes"`
	RuntimeSec  *int       `json:"runtime_sec"`
	LogPath     *string    `json:"log_path"`
	BackupPath  *string    `json:"backup_path"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	BackupJob   *BackupJob `json:"backup_job,omitempty"`
}

// NewClient creates a new API client
func NewClient() (*Client, error) {
	// Get API URL from environment or use default
	baseURL := getEnvWithDefault("API_URL", "https://api.dbbackup.io")

	// Get auth data from config
	token := config.GetToken()
	teamID := config.GetTeamId()

	if token == "" {
		return nil, fmt.Errorf("not authenticated. Please run 'dbbackup login' first")
	}

	if teamID == "" {
		return nil, fmt.Errorf("no team ID found. Please run 'dbbackup login' again")
	}

	client := &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Token:  token,
		TeamID: teamID,
	}

	logger.Debugf("Created API client for team %s with base URL: %s", teamID, baseURL)
	return client, nil
}

// makeRequest makes an authenticated HTTP request
func (c *Client) makeRequest(method, path string) (*http.Response, error) {
	url := c.BaseURL + path
	logger.Debugf("Making %s request to: %s", method, url)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication header
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	logger.Debugf("Response status: %d", resp.StatusCode)
	return resp, nil
}

// APIResponse represents a typical API response wrapper
type APIResponse struct {
	Data    []BackupJob `json:"data"`
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
}

// GetBackupJobs retrieves backup jobs for the team
func (c *Client) GetBackupJobs() ([]BackupJob, error) {
	path := fmt.Sprintf("/v1/teams/%s/backup_jobs", c.TeamID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// First, let's see what the actual response looks like
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	logger.Debugf("API response: %s", string(body))

	// Try to decode as wrapped response first
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Data != nil {
		logger.Debugf("Retrieved %d backup jobs from wrapped response", len(apiResp.Data))
		return apiResp.Data, nil
	}

	// If that fails, try to decode as direct array
	var jobs []BackupJob
	if err := json.Unmarshal(body, &jobs); err == nil {
		logger.Debugf("Retrieved %d backup jobs from direct array", len(jobs))
		return jobs, nil
	}

	// If both fail, return the original error with the response body for debugging
	return nil, fmt.Errorf("failed to decode response. Response body: %s", string(body))
}

// GetBackupJobDetails retrieves detailed information for a specific backup job
func (c *Client) GetBackupJobDetails(jobID int) (*BackupJob, error) {
	path := fmt.Sprintf("/v1/teams/%s/backup_jobs/%d", c.TeamID, jobID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	logger.Debugf("Job details API response: %s", string(body))

	// Try to decode as wrapped response first
	type JobDetailsResponse struct {
		Data   *BackupJob `json:"data"`
		Status string     `json:"status,omitempty"`
	}

	var detailsResp JobDetailsResponse
	if err := json.Unmarshal(body, &detailsResp); err == nil && detailsResp.Data != nil {
		logger.Debugf("Retrieved job details from wrapped response")
		return detailsResp.Data, nil
	}

	// If that fails, try to decode as direct object
	var job BackupJob
	if err := json.Unmarshal(body, &job); err == nil {
		logger.Debugf("Retrieved job details from direct object")
		return &job, nil
	}

	return nil, fmt.Errorf("failed to decode job details response. Response body: %s", string(body))
}

// RunBackupJob triggers a backup job to run
func (c *Client) RunBackupJob(jobID int) error {
	path := fmt.Sprintf("/v1/teams/%s/backup_jobs/%d/run", c.TeamID, jobID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to run backup job: %s", string(body))
	}

	logger.Debug("Backup job run triggered successfully")
	return nil
}

// GetServers retrieves all servers for the team
func (c *Client) GetServers() ([]Server, error) {
	path := fmt.Sprintf("/v1/teams/%s/servers", c.TeamID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	logger.Debugf("Servers API response: %s", string(body))

	// Try wrapped response first
	var apiResp struct {
		Data   []Server `json:"data"`
		Status string   `json:"status,omitempty"`
	}
	if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Data != nil {
		return apiResp.Data, nil
	}

	// Try direct array
	var servers []Server
	if err := json.Unmarshal(body, &servers); err == nil {
		return servers, nil
	}

	return nil, fmt.Errorf("failed to decode servers response. Response body: %s", string(body))
}

// GetServerDetails retrieves detailed information for a specific server
func (c *Client) GetServerDetails(serverID int) (*Server, error) {
	path := fmt.Sprintf("/v1/teams/%s/servers/%d", c.TeamID, serverID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try wrapped response first
	var detailsResp struct {
		Data   *Server `json:"data"`
		Status string  `json:"status,omitempty"`
	}
	if err := json.Unmarshal(body, &detailsResp); err == nil && detailsResp.Data != nil {
		return detailsResp.Data, nil
	}

	// Try direct object
	var server Server
	if err := json.Unmarshal(body, &server); err == nil {
		return &server, nil
	}

	return nil, fmt.Errorf("failed to decode server details response. Response body: %s", string(body))
}

// ValidateServer validates a server connection
func (c *Client) ValidateServer(serverID int) error {
	path := fmt.Sprintf("/v1/teams/%s/servers/%d/validate", c.TeamID, serverID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server validation failed: %s", string(body))
	}

	return nil
}

// GetDatabaseSources retrieves all database sources for the team
func (c *Client) GetDatabaseSources() ([]DatabaseSource, error) {
	path := fmt.Sprintf("/v1/teams/%s/database_sources", c.TeamID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try wrapped response first
	var apiResp struct {
		Data   []DatabaseSource `json:"data"`
		Status string           `json:"status,omitempty"`
	}
	if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Data != nil {
		return apiResp.Data, nil
	}

	// Try direct array
	var sources []DatabaseSource
	if err := json.Unmarshal(body, &sources); err == nil {
		return sources, nil
	}

	return nil, fmt.Errorf("failed to decode database sources response. Response body: %s", string(body))
}

// GetDatabaseSourceDetails retrieves detailed information for a specific database source
func (c *Client) GetDatabaseSourceDetails(sourceID int) (*DatabaseSource, error) {
	path := fmt.Sprintf("/v1/teams/%s/database_sources/%d", c.TeamID, sourceID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try wrapped response first
	var detailsResp struct {
		Data   *DatabaseSource `json:"data"`
		Status string          `json:"status,omitempty"`
	}
	if err := json.Unmarshal(body, &detailsResp); err == nil && detailsResp.Data != nil {
		return detailsResp.Data, nil
	}

	// Try direct object
	var source DatabaseSource
	if err := json.Unmarshal(body, &source); err == nil {
		return &source, nil
	}

	return nil, fmt.Errorf("failed to decode database source details response. Response body: %s", string(body))
}

// ValidateDatabaseSource validates a database source connection
func (c *Client) ValidateDatabaseSource(sourceID int) error {
	path := fmt.Sprintf("/v1/teams/%s/database_sources/%d/validate", c.TeamID, sourceID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("database source validation failed: %s", string(body))
	}

	return nil
}

// GetStorageDestinations retrieves all storage destinations for the team
func (c *Client) GetStorageDestinations() ([]StorageDestination, error) {
	path := fmt.Sprintf("/v1/teams/%s/storage_destinations", c.TeamID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try wrapped response first
	var apiResp struct {
		Data   []StorageDestination `json:"data"`
		Status string               `json:"status,omitempty"`
	}
	if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Data != nil {
		return apiResp.Data, nil
	}

	// Try direct array
	var destinations []StorageDestination
	if err := json.Unmarshal(body, &destinations); err == nil {
		return destinations, nil
	}

	return nil, fmt.Errorf("failed to decode storage destinations response. Response body: %s", string(body))
}

// GetStorageDestinationDetails retrieves detailed information for a specific storage destination
func (c *Client) GetStorageDestinationDetails(destID int) (*StorageDestination, error) {
	path := fmt.Sprintf("/v1/teams/%s/storage_destinations/%d", c.TeamID, destID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try wrapped response first
	var detailsResp struct {
		Data   *StorageDestination `json:"data"`
		Status string              `json:"status,omitempty"`
	}
	if err := json.Unmarshal(body, &detailsResp); err == nil && detailsResp.Data != nil {
		return detailsResp.Data, nil
	}

	// Try direct object
	var destination StorageDestination
	if err := json.Unmarshal(body, &destination); err == nil {
		return &destination, nil
	}

	return nil, fmt.Errorf("failed to decode storage destination details response. Response body: %s", string(body))
}

// ValidateStorageDestination validates a storage destination connection
func (c *Client) ValidateStorageDestination(destID int) error {
	path := fmt.Sprintf("/v1/teams/%s/storage_destinations/%d/validate", c.TeamID, destID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("storage destination validation failed: %s", string(body))
	}

	return nil
}

// GetBackupRunLogs retrieves logs for a specific backup run
func (c *Client) GetBackupRunLogs(runID int) (string, error) {
	path := fmt.Sprintf("/v1/teams/%s/backup_runs/%d/logs", c.TeamID, runID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// GetBackupRunsForJob retrieves backup runs for a specific backup job
func (c *Client) GetBackupRunsForJob(jobID int) ([]BackupRun, error) {
	// Use the correct endpoint for job runs
	path := fmt.Sprintf("/v1/teams/%s/backup_jobs/%d/backup_runs", c.TeamID, jobID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	logger.Debugf("Backup runs response: %s", string(body))

	// Try wrapped response first
	var apiResp struct {
		Data   []BackupRun `json:"data"`
		Status string      `json:"status,omitempty"`
	}
	if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Data != nil {
		return apiResp.Data, nil
	}

	// Try direct array
	var runs []BackupRun
	if err := json.Unmarshal(body, &runs); err == nil {
		return runs, nil
	}

	return nil, fmt.Errorf("failed to decode backup runs response. Response body: %s", string(body))
}

// GetBackupRunDetails retrieves detailed information for a specific backup run
func (c *Client) GetBackupRunDetails(runID int) (*BackupRun, error) {
	path := fmt.Sprintf("/v1/teams/%s/backup_runs/%d", c.TeamID, runID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try wrapped response first
	var detailsResp struct {
		Data   *BackupRun `json:"data"`
		Status string     `json:"status,omitempty"`
	}
	if err := json.Unmarshal(body, &detailsResp); err == nil && detailsResp.Data != nil {
		return detailsResp.Data, nil
	}

	// Try direct object
	var run BackupRun
	if err := json.Unmarshal(body, &run); err == nil {
		return &run, nil
	}

	return nil, fmt.Errorf("failed to decode backup run details response. Response body: %s", string(body))
}

// DownloadBackupRun downloads a backup run file
func (c *Client) DownloadBackupRun(runID int) (*http.Response, error) {
	path := fmt.Sprintf("/v1/teams/%s/backup_runs/%d/download", c.TeamID, runID)

	resp, err := c.makeRequest("GET", path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Return the response so caller can handle the download stream
	return resp, nil
}

// getEnvWithDefault returns the value of an environment variable or a default value if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
