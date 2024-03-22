package models

import (
	"net/http"
	"time"
)

type Middleware func(handler http.HandlerFunc) http.HandlerFunc

type Session struct {
	UserID         int       `json:"user_id"`
	ConnectionID   int       `json:"connection_id"`
	Username       string    `json:"username"`
	IpAddress      string    `json:"ip_address"`
	ExpirationTime time.Time `json:"expiration_time"`
}

type Credentials struct {
	Username string
	Password string
}

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

type MangaUser struct {
	Id              string `json:"id,omitempty"`
	LastChapterRead string `json:"last_chapter_read,omitempty"`
}

type TempUser struct {
	ConfirmID    string
	CreationTime time.Time
	User         User
}

type MailConfig struct {
	Email    string `json:"email_addr"`
	Auth     string `json:"email_auth"`
	Hostname string `json:"host"`
	Port     int    `json:"port"`
}
