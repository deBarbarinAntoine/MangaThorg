package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"mangathorg/internal/models"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type Logs []Log

type Log struct {
	Time       time.Time      `json:"time"`
	Level      string         `json:"level"`
	Message    string         `json:"message"`
	ReqId      int            `json:"req_id,omitempty"`
	User       models.Session `json:"user,omitempty"`
	ClientIP   string         `json:"client_ip,omitempty"`
	ReqMethod  string         `json:"req_method,omitempty"`
	ReqURL     string         `json:"req_url,omitempty"`
	HttpStatus int            `json:"http_status,omitempty"`
	ErrOutput  string         `json:"output,omitempty"`
}

var Logger *slog.Logger
var logs *os.File
var wg sync.WaitGroup

// setDailyTimer sets the first LogInit waiting time to match midnight.
func setDailyTimer() time.Duration {
	t := time.Now()
	t = t.Add(time.Hour * 24)
	n := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	d := n.Sub(t)
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}
	log.Println("setDailyTimer() value: ", d)
	return d
}

// LogInit is meant to be run as a goroutine to create a new log file every day
// appending the file's creation timestamp in its name.
func LogInit() {
	duration := setDailyTimer()
	var jsonHandler *slog.JSONHandler
	for {
		filename := "logs/logs_" + time.Now().Format(time.DateOnly) + ".log"
		var err error
		logs, err = os.Open(filename)
		if err != nil {
			var createErr error
			if logs != nil {
				logs.Close()
			}
			logs, createErr = os.Create(filename)
			if createErr != nil {
				log.Println(GetCurrentFuncName(), slog.Any("output", createErr))
			}
		}
		jsonHandler = slog.NewJSONHandler(logs, nil)
		Logger = slog.New(jsonHandler)
		time.Sleep(duration)
		duration = time.Hour * 24
	}
}

// fetchLogInfo retrieves all Log from `file` and stores it in *log.
func (log *Logs) fetchLogInfo(file string) {
	fmt.Println(GetCurrentFuncName())
	defer wg.Done()
	filename := "logs/" + file
	data, err := os.ReadFile(filename)
	if len(data) == 0 {
		return
	}
	lines := bytes.Split(data, []byte("\n"))
	var singleLog Log
	for _, line := range lines {
		err = json.Unmarshal(line, &singleLog)
		if err != nil {
			return
		}
		*log = append(*log, singleLog)
		fmt.Printf("singleLog: %#v\n", singleLog)
	}
	fmt.Printf("log: %#v\n", log)
}

func printFileNames(files []os.DirEntry) []string {
	var result []string
	for _, file := range files {
		result = append(result, file.Name())
	}
	return result
}

// RetrieveLogs fetches all Log from all files *.log in /logs directory
// and returns a Logs array.
func RetrieveLogs() (logArray Logs) {
	logFiles, err := os.ReadDir(Path + "logs/.")
	fmt.Printf("logFiles: %#v\n", printFileNames(logFiles))
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	} else {
		reg := regexp.MustCompile("^[a-zA-Z0-9-_]+\\.log$")
		for _, file := range logFiles {
			if reg.MatchString(file.Name()) {
				wg.Add(1)
				go logArray.fetchLogInfo(file.Name())
			}
		}
	}
	wg.Wait()
	logArray.sortLogs()
	return logArray
}

// sortLogs sort all Log from the newest to the oldest.
func (log *Logs) sortLogs() {
	sort.Slice(*log, func(i, j int) bool {
		return (*log)[i].Time.After((*log)[j].Time)
	})
}

// FetchLevelLogs filters Log returning only Log matching the given `level`.
func FetchAttrLogs(attr string, value string) Logs {
	attr = strings.ToLower(attr)
	logs := RetrieveLogs()
	var result Logs
	switch attr {
	case "level":
		switch strings.ToUpper(value) {
		case "INFO", "WARN", "ERROR":
			for _, singleLog := range logs {
				if singleLog.Level == strings.ToUpper(value) {
					result = append(result, singleLog)
				}
			}
			break
		default:
			return nil
		}
	case "user", "username":
		for _, singleLog := range logs {
			if strings.ToLower(singleLog.User.Username) == strings.ToLower(value) {
				result = append(result, singleLog)
			}
		}
		break
	default:
		return nil
	}
	result.sortLogs()
	return result
}
