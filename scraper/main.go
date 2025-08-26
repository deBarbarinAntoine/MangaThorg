package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
	
	"mangathorg/internal/models/api"
	"mangathorg/internal/utils"
)

const (
	baseURL = "https://api.mangadex.org"
	version = "v0.1.0"
)

var (
	// Client is the http.Client used for all API requests.
	Client = http.Client{
		Timeout: time.Second * 5,
	}
	
	NbRateLimitExceeded = 0
	
	ErrTooManyRequests = errors.New("too many requests")
)

type chapter struct {
	Title   string
	Volume  string
	Chapter string
	ID      string
	Scans   api.ApiChapterScan
}

func mangaFeedURL(id string) string {
	return fmt.Sprintf("%s/manga/%s/feed", baseURL, id)
}

func chapterURL(id string) string {
	return fmt.Sprintf("%s/at-home/server/%s", baseURL, id)
}

func scanImageURL(baseURL, hash, quality, filename string) string {
	return fmt.Sprintf("%s/%s/%s/%s", baseURL, quality, hash, filename)
}

func getVersion() string {
	return version
}

func usage() {
	fmt.Printf(":: MangaDex Scraper %s\n", getVersion())
	fmt.Println(":: This tool will download all chapters of a manga from MangaDex.")
	fmt.Println(":: Usage:")
	fmt.Println(":: >  Download manga with id <manga_id> to <destination> directory")
	fmt.Println(":: ==>    manga-scraper <destination> <manga_id>")
	fmt.Println(":: >  Display this message")
	fmt.Println(":: ==>    manga-scraper --help")
	fmt.Println(":: Credits:")
	fmt.Println(":: >  Thorgan (https://github.com/debarbarinantoine/)")
}

func main() {
	now := time.Now()
	helpCommands := []string{"--help", "-h", "-help", "help", "h"}
	versionCommands := []string{"--version", "-v", "-version", "version", "v"}
	if len(os.Args) == 2 {
		if slices.Contains(helpCommands, os.Args[1]) {
			usage()
			os.Exit(0)
		}
		if slices.Contains(versionCommands, os.Args[1]) {
			fmt.Printf(":: MangaDex Scraper %s\n", getVersion())
			os.Exit(0)
		}
	}
	if len(os.Args) < 3 {
		fmt.Println(":: [ERROR] Invalid arguments")
		usage()
		os.Exit(1)
	}
	
	destination := os.Args[1]
	mangaID := os.Args[2]
	feed, err := getMangaFeed(mangaID)
	if err != nil {
		ErrExit(err)
	}
	
	chapterMap := make(map[string]struct{}, len(feed.Data))
	chapters := make([]chapter, 0, len(feed.Data))
	for _, ch := range feed.Data {
		key := fmt.Sprintf("%s|%s", ch.Attributes.Volume, ch.Attributes.Chapter)
		if _, exists := chapterMap[key]; !exists {
			chapterMap[key] = struct{}{}
			currentChapter := chapter{
				Title:   strings.TrimSpace(ch.Attributes.Title),
				Volume:  strings.TrimSpace(ch.Attributes.Volume),
				Chapter: strings.TrimSpace(ch.Attributes.Chapter),
				ID:      strings.TrimSpace(ch.Id),
			}
			chapters = append(chapters, currentChapter)
		}
	}
	
	for i := 0; i < len(chapters); i++ {
		ch := chapters[i]
		fmt.Printf(":: [INFO] Getting scan hashes for Volume %s - Chapter %s - %s\n", ch.Volume, ch.Chapter, ch.Title)
		scan, err := getChapterScans(ch.ID)
		if err != nil {
			if errors.Is(err, ErrTooManyRequests) {
				fmt.Println()
				fmt.Println(":: [INFO] Rate limit exceeded")
				fmt.Println(":: [INFO] Waiting 10 seconds")
				for j := 0; j < 10; j++ {
					fmt.Print(".")
					time.Sleep(time.Second)
				}
				fmt.Println()
				i--
			}
			fmt.Printf(":: [ERROR] %s\n", err.Error())
			continue
		}
		chapters[i].Scans = *scan
		
		time.Sleep(time.Millisecond * 2000)
	}
	
	for i := 0; i < len(chapters); i++ {
		ch := chapters[i]
		
		var output string
		if ch.Title != "" {
			output = filepath.Join(destination, fmt.Sprintf("Volume %s", toFilename(ch.Volume)), fmt.Sprintf("Chapter %s - %s", toFilename(ch.Chapter), toFilename(ch.Title)))
		} else {
			output = filepath.Join(destination, fmt.Sprintf("Volume %s", toFilename(ch.Volume)), fmt.Sprintf("Chapter %s", toFilename(ch.Chapter)))
		}
		err = createDirs(output)
		if err != nil {
			ErrExit(err)
		}
		fmt.Printf(":: [INFO] Downloading Volume %s - Chapter %s - %s", ch.Volume, ch.Chapter, ch.Title)
		
		err = downloadChapter(output, ch)
		if err != nil {
			fmt.Println()
			if errors.Is(err, ErrTooManyRequests) {
				fmt.Println(":: [INFO] Rate limit exceeded")
				fmt.Println(":: [INFO] Waiting 10 seconds")
				for j := 0; j < 10; j++ {
					fmt.Print(".")
					time.Sleep(time.Second)
				}
				fmt.Println()
				i--
				continue
			}
			fmt.Printf(":: [ERROR] %s\n", err.Error())
			continue
		}
		
		time.Sleep(time.Millisecond * 2000)
	}
	
	fmt.Printf(":: [INFO] Download completed in %s\n", utils.DurationToString(time.Since(now)))
}

func toFilename(str string) string {
	regex := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	return regex.ReplaceAllString(str, "")
}

func getMangaFeed(id string) (*api.ApiMangaFeed, error) {
	feedURL := mangaFeedURL(id)
	var query = make(url.Values)
	query.Add("order[chapter]", "asc")
	query.Add("translatedLanguage[]", "en")
	query.Add("limit", "500")
	
	res, err := request(feedURL, query)
	if err != nil {
		return nil, err
	}
	
	feed := &api.ApiMangaFeed{}
	err = json.Unmarshal(res, feed)
	if err != nil {
		return nil, err
	}
	
	return feed, nil
}

func getChapterScans(id string) (*api.ApiChapterScan, error) {
	u := chapterURL(id)
	
	res, err := request(u, nil)
	if err != nil {
		return nil, err
	}
	
	ch := &api.ApiChapterScan{}
	err = json.Unmarshal(res, ch)
	if err != nil {
		return nil, err
	}
	
	if ch.Chapter.Data == nil && ch.Chapter.DataSaver == nil {
		return nil, fmt.Errorf("no scans found")
	}
	
	return ch, nil
}

func downloadChapter(destination string, ch chapter) error {
	var scanIDs []string
	var quality string
	if ch.Scans.Chapter.Data != nil {
		quality = "data"
		scanIDs = ch.Scans.Chapter.Data
	} else if ch.Scans.Chapter.DataSaver != nil {
		quality = "data-saver"
		scanIDs = ch.Scans.Chapter.DataSaver
	} else {
		return fmt.Errorf("no scans found")
	}
	
	for _, scanID := range scanIDs {
		err := downloadScanImage(ch.Scans.BaseUrl, ch.Scans.Chapter.Hash, quality, scanID, filepath.Join(destination, scanID))
		if err != nil {
			return err
		}
		fmt.Print(".")
		
		time.Sleep(time.Millisecond * 2000)
	}
	fmt.Println("[DONE]")
	
	return nil
}

func downloadScanImage(baseURL, hash, quality, filename, output string) error {
	u := scanImageURL(baseURL, hash, quality, filename)
	data, err := request(u, nil)
	if err != nil {
		return err
	}
	
	err = os.WriteFile(output, data, 0644)
	if err != nil {
		return err
	}
	
	return nil
}

func request(url string, query url.Values) ([]byte, error) {
	req, errReq := http.NewRequest(http.MethodGet, url, nil)
	if errReq != nil {
		return nil, errReq
	}
	
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	
	res, errRes := Client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	} else {
		return nil, errRes
	}
	
	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		return nil, errBody
	}
	
	var err error
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusTooManyRequests {
			err = ErrRateLimiter()
		} else {
			err = fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
	}
	
	return body, err
}

func createDirs(dirname string) error {
	return os.MkdirAll(dirname, 0755)
}

func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func ErrRateLimiter() error {
	NbRateLimitExceeded++
	if NbRateLimitExceeded > 5 {
		ErrExit(fmt.Errorf("rate limit exceeded %d times", NbRateLimitExceeded))
	}
	return ErrTooManyRequests
}

func ErrExit(err error) {
	fmt.Printf(":: [ERROR] %s\n", err.Error())
	fmt.Println(":: [WARN] Aborting operation...")
	os.Exit(1)
}
