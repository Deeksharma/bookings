package main

import (
	"github.com/Deeksharma/bookings/internal/config"
	"github.com/Deeksharma/bookings/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	//mux := pat.New()
	//
	//mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	//mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarter", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/reservation", handlers.Repo.SearchAvailability)
	mux.Post("/reservation", handlers.Repo.PostSearchAvailability)
	mux.Post("/reservation-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	// we'll create a file server, a place to go get static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(r chi.Router) {
		//r.Use(Auth)

		r.Get("/dashboard", handlers.Repo.AdminDashboard)
		r.Get("/reservationsnew", handlers.Repo.AdminNewReservation)
		r.Get("/reservationsall", handlers.Repo.AdminAllReservation)
		r.Get("/reservationscalendar", handlers.Repo.AdminReservationsCalender)

		r.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservation)
	})

	return mux
}
