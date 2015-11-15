package auth

import "os"

// Config base authentication configuration
type Config struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
}

// GitHubConfig github auth configuration
type GitHubConfig struct {
	*Config
}

// DefaultGitHubConfig default configuration used if not overridden
var DefaultGitHubConfig *GitHubConfig

func init() {
	DefaultGitHubConfig = &GitHubConfig{
		&Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			CallbackURL:  os.Getenv("GITHUB_CALLBACK_URL"),
		},
	}
}
