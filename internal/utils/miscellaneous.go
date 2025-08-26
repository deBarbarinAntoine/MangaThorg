package utils

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	// Path is the path to the project's root directory.
	Path = "./"
	
	// DataPath is the path to the cache directory.
	DataPath string = Path + "cache/"
)

var BaseURL = os.Getenv("BASE_URL")

// DurationToString -> just for fun ;)
func DurationToString(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%dh%0dm%0ds", hours, minutes, seconds)
	}
	return fmt.Sprintf("%dm%0ds", minutes, seconds)
}

// SetDailyTimer sets a waiting time to match a certain `hour`.
func SetDailyTimer(hour int) time.Duration {
	hour = hour % 24
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, t.Location())
	d := n.Sub(t)
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}
	log.Println("SetDailyTimer() value: ", DurationToString(d), "until", n.Format("02 Jan 15H04")) // verbose
	return d
}

// GetIP
//
//	@Description: gets the client's IP address according to the *http.Request.
//	@param r
//	@return string
func GetIP(r *http.Request) string {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")
	
	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String()
		}
	}
	
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}
	
	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1"
		}
		return ip
	}
	
	log.Fatalln(err)
	return ""
}

// GetCurrentFuncName
//
//	@Description: gets the name of the function that calls it.
//	@return string
func GetCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
