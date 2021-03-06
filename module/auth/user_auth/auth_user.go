package user_auth

import (
	"fmt"
	"net/http"

	"github.com/foxiswho/shop-go/middleware/session"
	"github.com/labstack/echo"
	"github.com/foxiswho/shop-go/models/auth"
)

const (
	DefaultKey  = "github.com/foxiswho/shop-go/modules/auth/user_auth"
	errorFormat = "[modules] ERROR! %s\n"
)

var (
	// RedirectUrl should be the relative URL for your login route
	RedirectUrl string = "/login"

	// RedirectParam is the query string parameter that will be set
	// with the page the user_service was trying to visit before they were
	// intercepted.
	RedirectParam string = "return_url"

	// SessionKey is the key containing the unique ID in your session
	SessionKey string = "AUTHUNIQUEID"
)

type AuthUser struct {
	auth.User
}

// shortcut to get AuthUser
func Default(c echo.Context) AuthUser {
	// return c.MustGet(DefaultKey).(auth)
	return c.Get(DefaultKey).(AuthUser)
}

// shortcut to get AuthUser
func DefaultGetUser(c echo.Context) auth.User {
	// return c.MustGet(DefaultKey).(auth)
	auth := c.Get(DefaultKey).(AuthUser)
	return auth.User
}

// AuthenticateSession will mark the session and user_service object as authenticated. Then
// the Login() user_service function will be called. This function should be called after
// you have validated a user_service.
func AuthenticateSession(s session.Session, user auth.User) error {
	user.Login()
	return UpdateUser(s, user)
}

func (a AuthUser) Logout(s session.Session) {
	Logout(s, a.User)
}

// Logout will clear out the session and call the Logout() user_service function.
func Logout(s session.Session, user auth.User) {
	user.Logout()
	s.Delete(SessionKey)
	s.Save()
}

// LoginRequired verifies that the current user_service is authenticated. Any routes that
// require a login should have this handler placed in the flow. If the user_service is not
// authenticated, they will be redirected to /login with the "next" get parameter
// set to the attempted URL.
func LoginRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			a := Default(c)
			if a.User.IsAuthenticated() == false {
				uri := c.Request().RequestURI
				path := fmt.Sprintf("%s?%s=%s", RedirectUrl, RedirectParam, uri)
				c.Redirect(http.StatusMovedPermanently, path)
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// UpdateUser updates the Admin object stored in the session. This is useful incase a change
// is made to the user_service model that needs to persist across requests.
func UpdateUser(s session.Session, user auth.User) error {
	s.Set(SessionKey, user.UniqueId())
	s.Save()
	return nil
}

func GetRoleId(c echo.Context) int {
	user := DefaultGetUser(c)
	return user.RoleId()
}

func GetRoleExtend(c echo.Context) []int {
	user := DefaultGetUser(c)
	return user.RoleExtend()
}
