package server

import (
	"log"
	"mangathorg/internal/api"
	"mangathorg/internal/utils"
	"mangathorg/router"
	"net/http"
)

// Run is the main function of the whole HTTP server:
// it initializes the routes, make the asset folder
// available to the clients, runs all needed goroutines
// and the ListenAndServe() function.
func Run() {
	// Initializing the routes
	router.Init()

	// Sending the assets to the clients
	fs := http.FileServer(http.Dir(utils.Path + "assets"))
	router.Mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Clear the cache before running the server
	api.EmptyCache()

	// Running the goroutine to change log file every given time
	go utils.LogInit()

	// Running the goroutine to automatically remove expired sessions every given time
	go utils.MonitorSessions()

	// Running the goroutine to automatically remove old TempUsers and LostUsers
	go utils.ManageTempUsers()

	// Running the goroutine to automatically remove old CacheData
	go api.CacheMonitor()

	// Delete all cache data on shutdown
	defer api.EmptyCache()

	// Running the server
	log.Fatalln(http.ListenAndServe(":8080", router.Mux))
}
