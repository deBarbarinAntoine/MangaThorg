package router

import (
	"mangathorg/controllers"
	"net/http"
)

var Mux = http.NewServeMux()

func Init() {
	Mux.HandleFunc("GET /{$}", controllers.IndexHandlerGetBundle)
	Mux.HandleFunc("PUT /{$}", controllers.IndexHandlerPutBundle)
	Mux.HandleFunc("DELETE /{$}", controllers.IndexHandlerDeleteBundle)
	Mux.HandleFunc("POST /add", controllers.IndexHandlerPutBundle)
	Mux.HandleFunc("GET /login", controllers.LoginHandlerGetBundle)
	Mux.HandleFunc("POST /login", controllers.LoginHandlerPostBundle)
	Mux.HandleFunc("GET /register", controllers.RegisterHandlerGetBundle)
	Mux.HandleFunc("POST /register", controllers.RegisterHandlerPostBundle)
	Mux.HandleFunc("GET /home", controllers.HomeHandlerGetBundle)
	Mux.HandleFunc("GET /logs", controllers.LogHandlerGetBundle)
	Mux.HandleFunc("GET /confirm", controllers.ConfirmHandlerGetBundle)
	Mux.HandleFunc("GET /logout", controllers.LogoutHandlerGetBundle)
	Mux.HandleFunc("GET /principal", controllers.PrincipalHandlerGetBundle)
	Mux.HandleFunc("GET /manga/{id}", controllers.MangaRequestHandlerGet)
	Mux.HandleFunc("GET /categories", controllers.TagsHandlerGetBundle)
	Mux.HandleFunc("GET /category/{tagId}", controllers.CategoryHandlerGetBundle)
	Mux.HandleFunc("GET /feed", controllers.FeedRequestHandlerGetBundle)
	Mux.HandleFunc("GET /chapter/{mangaId}/{chapterNb}/{chapterId}", controllers.ChapterHandlerGetBundle)
	Mux.HandleFunc("GET /mangatest", controllers.MangaWholeRequestHandlerGetBundle)
	Mux.HandleFunc("GET /covers/{manga}/{img}", controllers.CoversHandlerGetBundle)
	Mux.HandleFunc("GET /scan/{chapterId}/{quality}/{hash}/{img}", controllers.ScanHandlerGetBundle)

	// Handling MethodNotAllowed error on /
	Mux.HandleFunc("/{$}", controllers.IndexHandlerNoMethBundle)

	// Handling StatusNotFound error everywhere else
	Mux.HandleFunc("/", controllers.IndexHandlerOtherBundle)
}
