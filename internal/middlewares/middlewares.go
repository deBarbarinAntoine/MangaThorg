package middlewares

import (
	"errors"
	"log"
	"log/slog"
	"net/http"
	
	"mangathorg/internal/models/api"
	"mangathorg/internal/models/server"
	"mangathorg/internal/utils"
)

var LogId = 0

// Log is a models.Middleware that writes a series of information in logs/logs_<date>.log
// in JSON format: time, client's type, request Id (incremented int),
// user's models.Session (if logged), client IP, request Method, and request URL.
var Log server.Middleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		LogId++
		log.Println("middlewares.Log()")
		cookie, err := r.Cookie("session_id")
		if err != nil {
			utils.Logger.Info("Visitor", slog.Int("req_id", LogId), slog.String("client_ip", utils.GetIP(r)), slog.String("req_method", r.Method), slog.String("req_url", r.URL.String()))
		} else {
			utils.Logger.Info("User", slog.Int("req_id", LogId), slog.Any("user", utils.SessionsData[cookie.Value]), slog.String("client_ip", utils.GetIP(r)), slog.String("req_method", r.Method), slog.String("req_url", r.URL.String()))
		}
		next.ServeHTTP(w, r)
	}
}

// Guard is a models.Middleware that verify if a user has an opened session
// through the cookies and let it pass if ok, and redirects if not.
var Guard server.Middleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("middlewares.Guard()")
		
		// Checks if the user has a valid opened session
		ok := utils.CheckSession(r)
		if !ok {
			utils.Logger.Warn("Invalid session", slog.Int("req_id", LogId), slog.String("req_url", r.URL.String()), slog.Int("http_status", http.StatusUnauthorized))
			http.Redirect(w, r, "/login?err=restricted", http.StatusSeeOther)
			return
		}
		
		err := utils.RefreshSession(&w, r)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err), slog.Int("req_id", LogId))
		}
		
		// Use user data (e.g., display username) // testing
		// fmt.Fprintf(w, "Welcome, user %s", userData["user_id"])
		next.ServeHTTP(w, r)
	}
}

// SimpleGuard is a models.Middleware that verify if a user has an opened session
// through the cookies and let it pass if ok, and send an error if not.
var SimpleGuard server.Middleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("middlewares.SimpleGuard()")
		
		// Checks if the user has a valid opened session
		ok := utils.CheckSession(r)
		if !ok {
			utils.Logger.Warn("Invalid session", slog.Int("req_id", LogId), slog.String("req_url", r.URL.String()), slog.Int("http_status", http.StatusUnauthorized))
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}
		
		// retrieving the in-memory current session data
		cookie, err := r.Cookie("session_id")
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		
		// copying the cookie in the updatedCookie to retrieve it in the next handler
		cookie.Name = "updatedCookie"
		r.AddCookie(cookie)
		
		next.ServeHTTP(w, r)
	}
}

// UserCheck is a models.Middleware that checks if the client is logged,
// and if yes, it refreshes its sessionID
var UserCheck server.Middleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("middlewares.UserCheck()")
		exists := utils.CheckSession(r)
		if exists {
			err := utils.RefreshSession(&w, r)
			if err != nil {
				utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err), slog.Int("req_id", LogId))
			}
		}
		next.ServeHTTP(w, r)
	}
}

// OnlyVisitors is a models.Middleware that checks if the client is logged,
// and if yes, it redirects to the index handler.
var OnlyVisitors server.Middleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("middlewares.OnlyVisitors()")
		exists := utils.CheckSession(r)
		if exists {
			http.Redirect(w, r, "/principal", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// CheckApi is a models.Middleware that checks if the API is working properly,
// and if no, it redirects to the API error handler.
var CheckApi server.Middleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		_, err := api.Client.Head("https://api.mangadex.org/manga/tag")
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("an error occurred with the API")))
			http.Redirect(w, r, "/ErrorAPI", http.StatusSeeOther)
			return
		}
		
		next.ServeHTTP(w, r)
	}
}

// Join is used to concatenate various middlewares, for better visibility.
// it takes the http.HandlerFunc corresponding to the route, and then
// any number of models.Middleware that will be concatenated in order like this:
// middlewares[0](middlewares[1](middlewares[2](handlerFunc))).
func Join(handlerFunc http.HandlerFunc, middlewares ...server.Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handlerFunc = middlewares[i](handlerFunc)
	}
	return handlerFunc
}
