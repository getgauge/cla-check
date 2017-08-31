package configuration

import (
	"fmt"
	"log"
	"os"
)

// Dialect of the db to use for saving
func Dialect() string {
	return environmentVariable("DB_DIALECT")
}

// CookieName of the db to use for saving
func CookieName() string {
	return environmentVariable("COOKIE_NAME")
}

// Connection string to connect to preferred db
func Connection() string {
	return environmentVariable("DATABASE_URL")
}

// Port to use for the server
func Port() string {
	return environmentVariable("PORT")
}

//GithubKey to use for oAuth2
func GithubKey() string {
	return environmentVariable("GITHUB_KEY")
}

//GithubSecret to use for oAuth2
func GithubSecret() string {
	return environmentVariable("GITHUB_SECRET")
}

//GithubAuthCallback url for oAuth2
func GithubAuthCallback() string {
	return fmt.Sprintf("%s/auth/github/callback", environmentVariable("CALLBACK_HOST"))
}

//ContributorURL for heroku to list all contributors
func ContributorURL() string {
	return environmentVariable("CONTRIBUTOR_URL")
}

//AccessToken for creating comments on pr's and issues
func AccessToken() string {
	return environmentVariable("ACCESS_TOKEN")
}

func environmentVariable(variable string) string {
	value := os.Getenv(variable)
	if value == "" {
		log.Fatal(fmt.Sprintf("$%s must be set", variable))
	}
	return value
}
