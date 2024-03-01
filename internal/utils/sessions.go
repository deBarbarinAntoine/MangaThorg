package utils

import (
	"crypto/rand"
	"encoding/base64"
	"mangathorg/internal/models"
	"net/http"
	"time"
)

// In-memory Session data storage
var SessionsData = make(map[string]models.Session)

func retrieveSessions() []models.Session {
	var sessions []models.Session
	for _, session := range SessionsData {
		sessions = append(sessions, session)
	}
	return sessions
}

func newConnectionID() int {
	sessions := retrieveSessions()
	var id int
	var idFound bool
	for id = 1; !idFound; id++ {
		idFound = true
		for _, session := range sessions {
			if session.ConnectionID == id {
				idFound = false
			}
		}
	}
	id--
	return id
}

func OpenSession(w *http.ResponseWriter, username string, r *http.Request) {

	// Generate and set Session ID cookie
	sessionID := generateSessionID()
	// Generate expiration time for the cookie
	expirationTime := time.Now().Add(time.Minute)

	newCookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expirationTime,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(*w, newCookie)
	r.AddCookie(newCookie)

	user, _ := SelectUser(username)

	// Create Session data in memory
	SessionsData[sessionID] = models.Session{
		UserID:         user.Id,
		ConnectionID:   newConnectionID(),
		Username:       username,
		IpAddress:      GetIP(r),
		ExpirationTime: expirationTime,
	}
}

// CheckSession checks if there is a cookie in the request
// and if yes, it checks if the corresponding models.Session
// is still valid and returns true if all verifications are ok.
func CheckSession(r *http.Request) bool {
	// Extract session ID from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil || !validateSessionID(cookie.Value) {
		return false
	}
	// Retrieve user data from session
	session, ok := SessionsData[cookie.Value]
	if !ok {
		return ok
	}
	// Verify user IP address
	if session.IpAddress != GetIP(r) {
		return false
	}
	// Verify expiration time
	if session.ExpirationTime.Before(time.Now()) {
		return false
	}
	return true
}

func RefreshSession(w *http.ResponseWriter, r *http.Request) error {
	// generating new sessionID and new expiration time
	newSessionID := generateSessionID()
	newExpirationTime := time.Now().Add(time.Minute)

	var newCookie = &http.Cookie{
		Name:     "session_id",
		Value:    newSessionID,
		HttpOnly: true,
		Secure:   false, // Use only if using HTTPS
		Path:     "/",
		Expires:  newExpirationTime,
		SameSite: http.SameSiteStrictMode,
	}

	// setting the new cookie
	http.SetCookie(*w, newCookie)

	// retrieving the in-memory current session data
	cookie, err := r.Cookie("session_id")
	currentSessionData := SessionsData[cookie.Value]

	// updating the sessionID and expirationTime
	currentSessionData.ExpirationTime = newExpirationTime

	// deleting previous entry in the SessionsData map
	delete(SessionsData, cookie.Value)

	// setting the new entry in the SessionsData map
	SessionsData[newSessionID] = currentSessionData

	// adding the new cookie to the request to access it from the targeted handler with the Name "updatedCookie"
	newCookie.Name = "updatedCookie"
	r.AddCookie(newCookie)

	if err != nil {
		return err
	}
	return nil
}

func generateSessionID() string {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func validateSessionID(sessionID string) bool {
	_, ok := SessionsData[sessionID]
	return len(sessionID) == 88 && ok
}

func isExpired(session models.Session) bool {
	return session.ExpirationTime.Before(time.Now())
}

func cleanSessions() {
	for sessionID, session := range SessionsData {
		if isExpired(session) {
			delete(SessionsData, sessionID)
		}
	}
}

func MonitorSessions() {
	for {
		time.Sleep(time.Hour)
		cleanSessions()
	}
}
