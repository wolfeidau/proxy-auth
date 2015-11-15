package auth

import "encoding/json"

// GitHubUser represents a user in the github API
type GitHubUser struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// DecodeGitHubUser decode json into a user object
func DecodeGitHubUser(buf []byte) (*GitHubUser, error) {
	user := new(GitHubUser)

	err := json.Unmarshal(buf, user)

	return user, err
}
