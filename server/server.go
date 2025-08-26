package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	
	"mangathorg/internal/api"
	"mangathorg/internal/utils"
	"mangathorg/router"
)

// Run is the main function of the whole HTTP server:
// it initializes the routes, make the asset folder
// available to the clients, runs all needed goroutines
// and the ListenAndServe() function.
func Run() {
	
	// Initializing the port and the BaseURL
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	if utils.BaseURL == "" {
		utils.BaseURL = "http://localhost:" + port
	}
	
	// Initializing the routes
	router.Init()
	
	// Sending the assets to the clients
	fs := http.FileServer(http.Dir(utils.Path + "assets"))
	router.Mux.Handle("/static/", http.StripPrefix("/static/", fs))
	
	utils.InitUsers()
	
	// Running the goroutine to change log file every given time
	go utils.LogInit()
	
	// Running the goroutine to automatically remove expired sessions every given time
	go utils.MonitorSessions()
	
	// Running the goroutine to automatically remove old TempUsers and LostUsers
	go utils.ManageTempUsers()
	
	// Delete all cache data
	api.EmptyCache()
	
	// Running the goroutine to automatically remove old CacheData
	go api.CacheMonitor()
	
	// Waiting for the goroutines to be ready before starting the server
	time.Sleep(500 * time.Millisecond)
	
	// Logging the start of the server
	utils.Logger.Info(fmt.Sprintf("Server is listening on %s", utils.BaseURL))
	
	// Running the server
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", port), router.Mux))
}
