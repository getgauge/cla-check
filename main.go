package main

import (
	"errors"
	"net/http"
	"time"

	// optionally, used for session's encoder/decoder
	"github.com/getgauge/cla-check/comment"
	"github.com/getgauge/cla-check/configuration"
	"github.com/getgauge/cla-check/data"
	"github.com/gorilla/securecookie"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
)

var sessionsManager *sessions.Sessions

const (
	refererURL = "referer_url"
)

// These are some function helpers that you may use if you want
func init() {
	// attach a session manager
	cookieName := configuration.CookieName()
	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	hashKey := make([]byte, 32)
	blockKey := make([]byte, 32)
	secureCookie := securecookie.New(hashKey, blockKey)

	sessionsManager = sessions.New(sessions.Config{
		Cookie: cookieName,
		Encode: secureCookie.Encode,
		Decode: secureCookie.Decode,
	})
}

// GetProviderName is a function used to get the name of a provider
// for a given request. By default, this provider is fetched from
// the URL query string. If you provide it in a different way,
// assign your own function to this variable that returns the provider
// name for your request.
var GetProviderName = func(ctx context.Context) (string, error) {
	// try to get it from the url param "provider"
	if p := ctx.URLParam("provider"); p != "" {
		return p, nil
	}

	// try to get it from the url PATH parameter "{provider} or :provider or {provider:string} or {provider:alphabetical}"
	if p := ctx.Params().Get("provider"); p != "" {
		return p, nil
	}

	// try to get it from context's per-request storage
	if p := ctx.Values().GetString("provider"); p != "" {
		return p, nil
	}
	// if not found then return an empty string with the corresponding error
	return "", errors.New("you must select a provider")
}

/*
BeginAuthHandler is a convenience handler for starting the authentication process.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".

BeginAuthHandler will redirect the user to the appropriate authentication end-point
for the requested provider.

See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
func BeginAuthHandler(ctx context.Context) {
	url, err := GetAuthURL(ctx)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("%v", err)
		return
	}

	ctx.Redirect(url, iris.StatusTemporaryRedirect)
}

/*
GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.

It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider" or from the context's value of "provider" key.

I would recommend using the BeginAuthHandler instead of doing all of these steps
yourself, but that's entirely up to you.
*/
func GetAuthURL(ctx context.Context) (string, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return "", err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth(SetState(ctx))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}
	session := sessionsManager.Start(ctx)
	session.Set(providerName, sess.Marshal())
	return url, nil
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var SetState = func(ctx context.Context) string {
	state := ctx.URLParam("state")
	if len(state) > 0 {
		return state
	}

	return "state"

}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(ctx context.Context) string {
	return ctx.URLParam("state")
}

/*
CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all of the basic information about the user from the provider.

It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".

See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
var CompleteUserAuth = func(ctx context.Context) (goth.User, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}
	session := sessionsManager.Start(ctx)
	value := session.GetString(providerName)
	if value == "" {
		return goth.User{}, errors.New("session value for " + providerName + " not found")
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}

	user, err := provider.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, ctx.Request().URL.Query())
	if err != nil {
		return goth.User{}, err
	}

	session.Set(providerName, sess.Marshal())
	return provider.FetchUser(sess)
}

// Logout invalidates a user session.
func Logout(ctx context.Context) error {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return err
	}
	session := sessionsManager.Start(ctx)
	session.Delete(providerName)
	return nil
}

// Handlers

func contributorHandler(ctx context.Context) {
	nickName := ctx.URLParam("checkContributor")
	m := make(map[string]string, 0)

	if data.Signed(nickName) {
		m["username"] = nickName
		m["isContributor"] = "true"
	}

	ctx.JSON(m)
}

func contributorsHandler(ctx context.Context) {
	ctx.JSON(data.GetAll())
}

func logoutHandler(ctx context.Context) {
	Logout(ctx)
	ctx.Redirect("/", iris.StatusTemporaryRedirect)
}

func createTemplateData(u goth.User, refURL string) map[string]string {
	td := make(map[string]string, 0)
	td["Name"] = u.Name
	td["NickName"] = u.NickName
	td["Referer"] = refURL
	return td
}

func authCallbackHandler(ctx context.Context) {
	user, err := CompleteUserAuth(ctx)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Writef("%v", err)
		return
	}

	data.Save(data.User{
		Time:        time.Now().Format("_2 Jan 2006"),
		Name:        user.Name,
		Email:       user.Email,
		NickName:    user.NickName,
		UserID:      user.UserID,
		Description: user.Description,
	})
	ru := ctx.GetCookie(refererURL)
	td := createTemplateData(user, ru)
	ctx.ViewData("", td)
	comment.CreateRecheckComment(ru)
	if err := ctx.View("user.html"); err != nil {
		ctx.Writef("%v", err)
	}
}

func providerHandler(ctx context.Context) {
	// try to get the user without re-authenticating
	if u, err := CompleteUserAuth(ctx); err == nil {
		td := createTemplateData(u, ctx.GetCookie(refererURL))
		ctx.ViewData("", td)
		if err := ctx.View("user.html"); err != nil {
			ctx.Writef("%v", err)
		}
	} else {
		BeginAuthHandler(ctx)
	}
}

func setReferer(ctx context.Context) {
	c := http.Cookie{
		Name:  refererURL,
		Value: ctx.Request().Referer(),
	}
	ctx.SetCookie(&c)
}

func defaultHandler(ctx context.Context) {
	if ctx.Request().Referer() != "" {
		setReferer(ctx)
	}
	if err := ctx.View("cla.html"); err != nil {
		ctx.Writef("%v", err)
	}
}

func main() {
	port := configuration.Port()

	db := data.Init()
	defer db.Close()

	goth.UseProviders(
		github.New(configuration.GithubKey(), configuration.GithubSecret(), configuration.GithubAuthCallback()),
	)

	app := iris.New()

	app.StaticWeb("/static", "./resources")
	app.RegisterView(iris.HTML("./templates", ".html").Layout("layout.html"))

	app.Get("/", defaultHandler)
	app.Get("/auth/{provider}", providerHandler)
	app.Get("/auth/{provider}/callback", authCallbackHandler)
	app.Get("/logout/{provider}", logoutHandler)
	app.Get("/contributor", contributorHandler)
	app.Get(configuration.ContributorURL(), contributorsHandler)
	app.Run(iris.Addr(":" + port))
}
