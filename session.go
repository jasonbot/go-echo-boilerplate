package quoteapi

import (
	"errors"
)

type sessionData struct {
	SessionID string `json:"session_id"`
	Username  string `json:"username"`
}

type sessionManager struct {
	sessionStore Datastore
}

func (manager *sessionManager) SignUp(username string, password string) (UserSession, error) {
	return nil, errors.New("Unimplemented")
}

func (manager *sessionManager) SignIn(username string, password string) (UserSession, error) {
	return nil, errors.New("Unimplemented")
}

func (manager *sessionManager) GetSession(sessionID string) (UserSession, error) {
	return nil, errors.New("Unimplemented")
}

// UserLogin handles logging a user in by session
type UserLogin interface {
	SignUp(username string, password string) (UserSession, error)
	SignIn(username string, password string) (UserSession, error)
	GetSession(sessionID string) (UserSession, error)
}

// UserSession represents the user's currenly logged in session
type UserSession interface {
	User() User
	SignOut()
}

// GetUserLogin Creates a session store based on its underlying data store
func GetUserLogin(dataStore Datastore) (UserLogin, error) {
	if dataStore != nil {
		return &sessionManager{sessionStore: dataStore}, nil
	}
	return nil, errors.New("No data store specified")
}
