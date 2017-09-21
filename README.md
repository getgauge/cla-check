# cla-check
Cla-check is an application which serves `Gauge` CLA page. With the integration of [cla-bot](https://github.com/ColinEberhardt/cla-bot) it checks if a contributors of a Pull Request have signed the CLA or not.
If all the commiters in PR have signed the CLA it adds `cla-signed` label to PR.

## Running Locally

### Prerequisites

* golang
* govender

### steps

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

