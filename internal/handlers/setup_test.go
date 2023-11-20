package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/Deeksharma/bookings/internal/config"
	"github.com/Deeksharma/bookings/internal/models"
	"github.com/Deeksharma/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var app config.AppConfig
var sessionManager *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	// change it to true when in production

	app.InProduction = false
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true // cookie will persist even after the browser is closed
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction // since localhost is not secure

	mailChan := make(chan models.MailData)
	app.MailChannel = mailChan
	defer close(app.MailChannel)

	listenForMail()

	cache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create templates: ", err)
	}
	app.SessionManager = sessionManager
	app.TemplateCache = cache
	app.UseCache = true // this needs to be true in case of tests otherwise render.CreateTemplate will be called again
	render.NewRenderer(&app)

	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	repo := NewTestingRepo(&app) // use our repo model to create a test db repo
	NewHandlers(repo)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarter", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/contact", Repo.Contact)

	mux.Get("/reservation", Repo.SearchAvailability)
	mux.Post("/reservation", Repo.PostSearchAvailability)
	mux.Post("/reservation-json", Repo.AvailabilityJSON)

	mux.Get("/make-reservation", Repo.MakeReservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	// we'll create a file server, a place to go get static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
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

// CreateTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) { // find the template and the layouts associated with it
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		} else {
			fmt.Println("no matches found!")
		}
		myCache[name] = templateSet
	}
	return myCache, nil
}

func listenForMail() {
	go func() {
		for{
			_ = <- app.MailChannel
		}
	}()
}