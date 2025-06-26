package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configDir  = ".dbbackup"
	configFile = "config"
)

// Init initializes the Viper configuration
func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDirPath := filepath.Join(home, configDir)

	// Ensure config directory exists
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set up Viper
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDirPath)

	// Set default values
	viper.SetDefault("token", "")
	viper.SetDefault("teamId", "")

	// Try to read config file (it's okay if it doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

// SetToken updates the token in the configuration
func SetToken(token string) error {
	if err := Init(); err != nil {
		return err
	}

	viper.Set("token", token)

	// Try WriteConfig first, if it fails try SafeWriteConfig
	if err := viper.WriteConfig(); err != nil {
		return viper.SafeWriteConfig()
	}
	return nil
}

// SetTeamId updates the teamId in the configuration
func SetTeamId(teamId string) error {
	if err := Init(); err != nil {
		return err
	}

	viper.Set("teamId", teamId)

	// Try WriteConfig first, if it fails try SafeWriteConfig
	if err := viper.WriteConfig(); err != nil {
		return viper.SafeWriteConfig()
	}
	return nil
}

// SetAuthData updates both token and teamId in the configuration
func SetAuthData(token, teamId string) error {
	if err := Init(); err != nil {
		return err
	}

	viper.Set("token", token)
	viper.Set("teamId", teamId)

	// Try WriteConfig first, if it fails try SafeWriteConfig
	if err := viper.WriteConfig(); err != nil {
		// If config file doesn't exist, create it
		return viper.SafeWriteConfig()
	}
	return nil
}

// GetToken returns the current token
func GetToken() string {
	if err := Init(); err != nil {
		return ""
	}
	return viper.GetString("token")
}

// GetTeamId returns the current teamId
func GetTeamId() string {
	if err := Init(); err != nil {
		return ""
	}
	return viper.GetString("teamId")
}

// IsAuthenticated checks if the user is authenticated (has a token)
func IsAuthenticated() bool {
	token := GetToken()
	return token != ""
}

// ClearToken removes the authentication token
func ClearToken() error {
	if err := Init(); err != nil {
		return err
	}

	viper.Set("token", "")

	// Try WriteConfig first, if it fails try SafeWriteConfig
	if err := viper.WriteConfig(); err != nil {
		return viper.SafeWriteConfig()
	}
	return nil
}

// ClearAuthData removes both token and teamId
func ClearAuthData() error {
	if err := Init(); err != nil {
		return err
	}

	viper.Set("token", "")
	viper.Set("teamId", "")

	// Try WriteConfig first, if it fails try SafeWriteConfig
	if err := viper.WriteConfig(); err != nil {
		return viper.SafeWriteConfig()
	}
	return nil
}
