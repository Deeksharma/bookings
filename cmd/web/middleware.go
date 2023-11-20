package main

import (
	"fmt"
	"github.com/Deeksharma/bookings/internal/helpers"
	"github.com/justinas/nosurf"
	"net/http"
)

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

//NoSurf adds CSRF protection to all post requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		//Name:       "",
		//Value:      "",
		Path: "/",
		//Domain:     "",
		//Expires:    time.Time{},
		//RawExpires: "",
		//MaxAge:     0,
		Secure:   app.InProduction,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		//Raw:        "",
		//Unparsed:   nil,
	})
	return csrfHandler
}

// SessionLoad loads and save session at every request
func SessionLoad(next http.Handler) http.Handler {
	return sessionManager.LoadAndSave(next)
}

// Auth protects root in our root file, only people logged in will be able to access further
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) { // typecasting this function as http.HandlerFunc
		if !helpers.IsAuthenticated(request) {
			sessionManager.Put(request.Context(), "error", "Login first!!")
			http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
