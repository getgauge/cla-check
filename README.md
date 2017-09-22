# CLA checker
`cla-check` promtps a contributor to sign the Gauge Contributor License agreement. It is an end-point to the excellent [cla-bot](https://github.com/ColinEberhardt/cla-bot). It's triggered on PR submission label's and adds a `cla-signed` label for authors who've accepted the CLA. 

## Running Locally

### Pre-Requisites

* golang
* govender

### Steps

* clone this repository
* run test by executing `go test ./...`
* set following required env variables
  * COOKIE_NAME :- the cookie which you wnat to use for auth 
  * PORT :- port where you wnat to run your app
  * DB_DIALECT :- the database which you want to use
  * DATABASE_URL :- URL of the database
  * GITHUB_KEY :- github client ID for OAuth app
  * GITHUB_SECRET :- gitub client secred for your OAuth app
  * CALLBACK_HOST :- A callback URL for OAuth to redirect afetr authentication
  * CONTRIBUTOR_URL :- a secret url to list all users from database
  * ACCESS_TOKEN :- A personal github access token to interact with github PRs
* Run the app by eecuting `go run main.go`

Deploy this to heroku if you want your own instance.
