package quoteapi

import (
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type sessionData struct {
	SessionID string `json:"session_id"`
	Username  string `json:"username"`
}

type sessionManager struct {
	sessionStore Datastore
}

type localSession struct {
	sessionStore    Datastore
	sessionData     sessionData
	sessionUserData userData
}

func (session *localSession) User() User {
	return &session.sessionUserData
}

func (session *localSession) SessionID() string {
	return session.sessionData.SessionID
}

func (session *localSession) SignOut() error {
	return session.sessionStore.DeleteRecord(
		"session",
		session.sessionData.SessionID,
	)
}

func (manager *sessionManager) SignUp(username string, password string) (UserSession, error) {
	var user userData
	if err := manager.sessionStore.LoadRecord("user", &user, username); err == nil {
		return nil, errors.New("User exists")
	}

	user.Username = username
	passwordString, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user.Password = base64.StdEncoding.EncodeToString(passwordString)
	user.PopulateFields()
	manager.sessionStore.SaveRecord("user", user, user.Username)

	return manager.SignIn(username, password)
}

func (manager *sessionManager) SignIn(username string, password string) (UserSession, error) {
	var user userData
	if err := manager.sessionStore.LoadRecord("user", &user, username); err != nil {
		return nil, err
	}

	passwordBytes, err := base64.RawStdEncoding.DecodeString(user.Password)

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(passwordBytes, []byte(password)); err != nil {
		return nil, err
	}

	newSessionID := uuid.New().String()
	newSession := sessionData{SessionID: newSessionID, Username: user.Username}
	manager.sessionStore.SaveRecord("session", newSession, newSession.SessionID)

	return &localSession{
		sessionStore:    manager.sessionStore,
		sessionData:     newSession,
		sessionUserData: user}, nil
}

func (manager *sessionManager) GetSession(sessionID string) (UserSession, error) {
	var userSession sessionData
	if err := manager.sessionStore.LoadRecord("session", userSession, sessionID); err != nil {
		return nil, errors.New("Session gone")
	}

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
	SessionID() string
	SignOut() error
}

// GetUserLogin Creates a session store based on its underlying data store
func GetUserLogin(dataStore Datastore) (UserLogin, error) {
	if dataStore != nil {
		return &sessionManager{sessionStore: dataStore}, nil
	}
	return nil, errors.New("No data store specified")
}
