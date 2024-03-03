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
	Id             int       `json:"id"`
	CreationTime   time.Time `json:"creation_time"`
	LastConnection time.Time `json:"last_connection"`
	Username       string    `json:"username"`
	HashedPwd      string    `json:"hash"`
	Salt           string    `json:"salt"`
	Email          string    `json:"email"`
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

type ApiManga struct {
	Result string `json:"result"`
	Errors []struct {
		Id      string      `json:"id"`
		Status  int         `json:"status"`
		Title   string      `json:"title"`
		Detail  string      `json:"detail"`
		Context interface{} `json:"context"`
	} `json:"errors,omitempty"`
	Response string `json:"response,omitempty"`
	Data     []struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Title struct {
				En string `json:"en"`
			} `json:"title"`
			AltTitles []struct {
				Zh   string `json:"zh,omitempty"`
				Ko   string `json:"ko,omitempty"`
				En   string `json:"en,omitempty"`
				Ja   string `json:"ja,omitempty"`
				De   string `json:"de,omitempty"`
				JaRo string `json:"ja-ro,omitempty"`
				Ru   string `json:"ru,omitempty"`
				Ar   string `json:"ar,omitempty"`
				Uk   string `json:"uk,omitempty"`
				KoRo string `json:"ko-ro,omitempty"`
				Mn   string `json:"mn,omitempty"`
				PtBr string `json:"pt-br,omitempty"`
			} `json:"altTitles"`
			Description struct {
				En   string `json:"en"`
				Uk   string `json:"uk,omitempty"`
				Ko   string `json:"ko,omitempty"`
				Ru   string `json:"ru,omitempty"`
				Zh   string `json:"zh,omitempty"`
				PtBr string `json:"pt-br,omitempty"`
			} `json:"description"`
			IsLocked bool `json:"isLocked"`
			Links    struct {
				Al    string `json:"al,omitempty"`
				Ap    string `json:"ap,omitempty"`
				Mu    string `json:"mu"`
				Raw   string `json:"raw,omitempty"`
				Bw    string `json:"bw,omitempty"`
				Kt    string `json:"kt,omitempty"`
				Amz   string `json:"amz,omitempty"`
				Ebj   string `json:"ebj,omitempty"`
				Mal   string `json:"mal,omitempty"`
				Nu    string `json:"nu,omitempty"`
				Engtl string `json:"engtl,omitempty"`
				Cdj   string `json:"cdj,omitempty"`
			} `json:"links"`
			OriginalLanguage       string  `json:"originalLanguage"`
			LastVolume             string  `json:"lastVolume"`
			LastChapter            string  `json:"lastChapter"`
			PublicationDemographic *string `json:"publicationDemographic"`
			Status                 string  `json:"status"`
			Year                   *int    `json:"year"`
			ContentRating          string  `json:"contentRating"`
			Tags                   []struct {
				Id         string `json:"id"`
				Type       string `json:"type"`
				Attributes struct {
					Name struct {
						En string `json:"en"`
					} `json:"name"`
					Description struct {
					} `json:"description"`
					Group   string `json:"group"`
					Version int    `json:"version"`
				} `json:"attributes"`
				Relationships []interface{} `json:"relationships"`
			} `json:"tags"`
			State                          string    `json:"state"`
			ChapterNumbersResetOnNewVolume bool      `json:"chapterNumbersResetOnNewVolume"`
			CreatedAt                      time.Time `json:"createdAt"`
			UpdatedAt                      time.Time `json:"updatedAt"`
			Version                        int       `json:"version"`
			AvailableTranslatedLanguages   []string  `json:"availableTranslatedLanguages"`
			LatestUploadedChapter          string    `json:"latestUploadedChapter"`
		} `json:"attributes"`
		Relationships []struct {
			Id      string `json:"id"`
			Type    string `json:"type"`
			Related string `json:"related,omitempty"`
		} `json:"relationships"`
	} `json:"data,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
	Total  int `json:"total,omitempty"`
}
