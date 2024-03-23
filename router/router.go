package router

import (
	"mangathorg/controllers"
	"net/http"
)

// Mux is the server's ServeMux.
var Mux = http.NewServeMux()

// Init initializes all routes.
func Init() {
	Mux.HandleFunc("GET /{$}", controllers.RootHandlerGetBundle)
	Mux.HandleFunc("GET /login", controllers.LoginHandlerGetBundle)
	Mux.HandleFunc("POST /login", controllers.LoginHandlerPostBundle)
	Mux.HandleFunc("GET /register", controllers.RegisterHandlerGetBundle)
	Mux.HandleFunc("POST /register", controllers.RegisterHandlerPostBundle)
	Mux.HandleFunc("GET /forgot-password", controllers.ForgotPasswordHandlerGetBundle)
	Mux.HandleFunc("POST /forgot-password", controllers.ForgotPasswordHandlerPostBundle)
	Mux.HandleFunc("GET /update-credentials", controllers.UpdateCredentialsHandlerGetBundle)
	Mux.HandleFunc("POST /update-credentials/{id}", controllers.UpdateCredentialsHandlerPostBundle)
	Mux.HandleFunc("GET /profile", controllers.ProfileHandlerGetBundle)
	Mux.HandleFunc("POST /profile", controllers.ProfileHandlerPostBundle)
	Mux.HandleFunc("GET /home", controllers.HomeHandlerGetBundle)
	Mux.HandleFunc("GET /confirm", controllers.ConfirmHandlerGetBundle)
	Mux.HandleFunc("GET /logout", controllers.LogoutHandlerGetBundle)
	Mux.HandleFunc("GET /principal", controllers.PrincipalHandlerGetBundle)
	Mux.HandleFunc("GET /manga/{id}", controllers.MangaRequestHandlerGet)
	Mux.HandleFunc("GET /categories", controllers.TagsHandlerGetBundle)
	Mux.HandleFunc("GET /category/{tagId}", controllers.CategoryHandlerGetBundle)
	Mux.HandleFunc("GET /category/{group}/{name}", controllers.CategoryNameHandlerGetBundle)
	Mux.HandleFunc("GET /search", controllers.SearchHandlerGetBundle)
	Mux.HandleFunc("GET /chapter/{mangaId}/{chapterNb}/{chapterId}", controllers.ChapterHandlerGetBundle)
	Mux.HandleFunc("GET /covers/{manga}/{img}", controllers.CoversHandlerGetBundle)
	Mux.HandleFunc("GET /scan/{chapterId}/{quality}/{hash}/{img}", controllers.ScanHandlerGetBundle)
	Mux.HandleFunc("POST /favorite/{mangaId}", controllers.FavoriteHandlerPostBundle)
	Mux.HandleFunc("DELETE /favorite/{mangaId}", controllers.FavoriteHandlerDeleteBundle)
	Mux.HandleFunc("PUT /favorite/{mangaId}", controllers.BannerHandlerPutBundle)

	// !! TESTING: this route is only for testing purposes for now. You need to disable it if you want to deploy the server.
	Mux.HandleFunc("GET /logs", controllers.LogHandlerGetBundle)

	// Handling StatusNotFound error everywhere else
	Mux.HandleFunc("/", controllers.ErrorHandlerBundle)
}
