package auth

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
var defaultGitHubConfig *GitHubConfig

// SetGitHubConfig Override the GitHub oauth configuration
func SetGitHubConfig(c *GitHubConfig) {
	defaultGitHubConfig = c
}
