package api

import (
	"time"
)

// Status is an enum-like variable for status keys.
var Status = struct {
	LastUploaded string
	Popular      string
	Tags         string
	Categories   string
	Mangas       string
	MangaFeeds   string
	ChaptersScan string
	MangaStats   string
}{
	LastUploaded: "last_uploaded",
	Popular:      "popular",
	Tags:         "tags",
	Categories:   "categories",
	Mangas:       "mangas",
	MangaFeeds:   "manga_feeds",
	ChaptersScan: "chapters_scan",
	MangaStats:   "manga_stats",
}

// StatusCache is the data structure of the status.json cache file.
type StatusCache struct {
	LastUploaded bool     `json:"last_uploaded"`
	Popular      bool     `json:"popular"`
	Tags         bool     `json:"tags"`
	Categories   []string `json:"categories"`
	Mangas       []string `json:"mangas"`
	MangaFeeds   []string `json:"manga_feeds"`
	ChaptersScan []string `json:"chapters_scan"`
	MangaStats   []string `json:"manga_stats"`
}

// CacheData is the data structure for all data stored in the cache.
type CacheData []SingleCacheData

// SingleCacheData is the data structure for every single data record stored in
// the cache.
type SingleCacheData struct {
	Id          string      `json:"id"`
	UpdatedTime time.Time   `json:"updated_time"`
	Order       string      `json:"order"`
	Offset      int         `json:"offset"`
	Data        interface{} `json:"data"`
}
