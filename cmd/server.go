package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	quoteapi "github.com/jasonbot/blueowl"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type locationResult struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func getLocation(originIP string) string {
	remoteURL := fmt.Sprintf("http://ip-api.com/json/%s", url.PathEscape(originIP))

	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, remoteURL, nil)
	if err != nil {
		return "Not on the internet"
	}

	res, err := client.Do(req)
	if err != nil {
		return "Nowhere"
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "???????"
	}

	location := locationResult{
		City:    "San Francisco",
		Country: "US",
	}
	if err := json.Unmarshal(body, &location); err != nil {
		return "???"
	}

	return fmt.Sprintf("%s, %s", location.City, location.Country)
}

func sessionMiddleware(datastore quoteapi.Datastore) echo.MiddlewareFunc {
	sessionStore, err := quoteapi.GetUserLogin(datastore)

	if err != nil {
		panic("Session store failed")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("session-store", sessionStore)

			cookie, err := c.Cookie("session-id")

			c.Logger().Debugf("Cookie fetch: %v %v", cookie, err)

			if err != nil && cookie != nil && cookie.Value != "" {
				session, err := sessionStore.GetSession(cookie.Value)

				if err != nil {
					c.Logger().Debugf("Failed to fetch session store: %v", err)
				}

				if session != nil {
					c.Logger().Debugf("User session found: %v", session)
					c.Set("user-session", session)
				}
			}

			if err := next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}

type usernamePassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Debug = true

	datastore, err := quoteapi.GetBoltStore("./prod_api")
	if err != nil {
		panic("Can't open datastore")
	}
	defer datastore.Close()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(sessionMiddleware(datastore))

	e.GET("/", hello)

	e.POST("/signup", func(c echo.Context) error {
		c.Logger().Debug("Sign up")
		var loginInfo usernamePassword
		c.Bind(&loginInfo)

		sessionStore, ok := c.Get("session-store").(quoteapi.UserLogin)

		c.Logger().Debugf("Session store: %v", sessionStore)

		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "No session store")
		}

		location := getLocation(c.RealIP())

		session, err := sessionStore.SignUp(loginInfo.Username, loginInfo.Password, location)

		if err != nil {
			c.Logger().Debugf("Error signing up: %v %v", loginInfo, err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		cookie := http.Cookie{
			Name:    "session-id",
			Value:   session.SessionID(),
			Expires: time.Now().Add(24 * time.Hour),
		}
		c.SetCookie(&cookie)

		return c.JSON(http.StatusOK, session.User())
	})

	e.POST("/logout", func(c echo.Context) error {
		sessionStore, ok := c.Get("user-session").(quoteapi.UserSession)

		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Not logged in")
		}

		if err := sessionStore.SignOut(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.String(http.StatusOK, "")
	})

	e.POST("/login", func(c echo.Context) error {
		var loginInfo usernamePassword
		c.Bind(&loginInfo)

		sessionStore, ok := c.Get("session-store").(quoteapi.UserLogin)

		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "No session store")
		}

		location := getLocation(c.RealIP())

		session, err := sessionStore.SignIn(loginInfo.Username, loginInfo.Password, location)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		cookie := http.Cookie{
			Name:    "session-id",
			Value:   session.SessionID(),
			Expires: time.Now().Add(24 * time.Hour),
		}
		c.SetCookie(&cookie)

		return c.JSON(http.StatusOK, session.User())
	})

	e.Logger.Fatal(e.Start(":8000"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
