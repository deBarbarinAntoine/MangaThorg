package models

import (
	"net/http"
	"time"
)

// Middleware is the type used for all middlewares.
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// Session is the structure used for any user session.
type Session struct {
	UserID         int       `json:"user_id"`
	ConnectionID   int       `json:"connection_id"`
	Username       string    `json:"username"`
	IpAddress      string    `json:"ip_address"`
	ExpirationTime time.Time `json:"expiration_time"`
}

// Credentials is the structure used to authenticate a user at login.
type Credentials struct {
	Username string
	Password string
}

// User is the structure used to store all user related data.
type User struct {
	Id             int         `json:"id"`
	CreationTime   time.Time   `json:"creation_time"`
	LastConnection time.Time   `json:"last_connection"`
	Username       string      `json:"username"`
	Avatar         string      `json:"avatar,omitempty"`
	HashedPwd      string      `json:"hash"`
	Salt           string      `json:"salt"`
	Email          string      `json:"email"`
	MangaBanner    MangaUser   `json:"manga_banner"`
	Favorites      []MangaUser `json:"favorites"`
}

// MangaUser is the structure used for all user related mangas.
type MangaUser struct {
	Id              string `json:"id,omitempty"`
	LastChapterRead string `json:"last_chapter_read,omitempty"`
}

// TempUser is the structure for any temporary user (waiting to be confirmed or
// which password has been forgotten).
type TempUser struct {
	ConfirmID    string
	CreationTime time.Time
	User         User
}

// MailConfig is the structure used to retrieve the sending mail's configuration.
type MailConfig struct {
	Email    string `json:"email_addr"`
	Auth     string `json:"email_auth"`
	Hostname string `json:"host"`
	Port     int    `json:"port"`
}
