package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/Deeksharma/bookings/internal/config"
	"github.com/Deeksharma/bookings/internal/driver"
	"github.com/Deeksharma/bookings/internal/handlers"
	"github.com/Deeksharma/bookings/internal/helpers"
	"github.com/Deeksharma/bookings/internal/models"
	"github.com/Deeksharma/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var sessionManager *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	//http.HandleFunc("/divide", handlers.Repo.Divide)
	db, err := Run() // we should pass our db object to handlers so that they can interact with it
	if err != nil {
		log.Fatal("Cannot start application, error:", err)
	}
	defer func() {
		db.SQL.Close()
		close(app.MailChannel)
	}()

	fmt.Println("Starting mail listener...")
	listenForMail() // listen for mails

	// using inbuilt go pkg
	/*from := "me@gmail.com"
	//to := "deeksha.sharma@zomato.com"

	// acc to the functionality of standard library we must have some means of authenticating mail server
	// i.e, give your credentials to send mail

	auth := smtp.PlainAuth("", from, "", "localhost")

	err = smtp.SendMail("localhost:1025", auth, from, []string{"you@here.com"}, []byte("Hello"))
	if err != nil {
		log.Println(err)
	}*/

	fmt.Printf("Starting application on port %v", portNumber)
	//http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// just leave the main function alone and test other things
func Run() (*driver.DB, error) {
	// what am I going to put in session, primitives are already registered
	//now we need to register type models.Reservation so that we can store value in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})

	mailChan := make(chan models.MailData)
	app.MailChannel = mailChan // make a listener for this

	// change it to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true // cookie will persist even after the browser is closed
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction // since localhost is not secure

	// connect to DB
	log.Println("Connecting to database...")
	ctx := context.Background()
	db, err := driver.ConnectSQL(ctx, "host=localhost port=5432 dbname=bookings user=deekshasharma password=")
	if err != nil {
		log.Fatal("Cannot connect to database. Dying!!!")
		return nil, err
	}
	log.Println("Connected to database!")
	// defer db.SQL.Close() // should not be put here because it'll start the pool and closes it in this function only

	cache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create templates: ", err)
		return nil, err
	}
	app.SessionManager = sessionManager
	app.TemplateCache = cache
	app.UseCache = false // runtime changes to html file wont reflect since the templates are rendered once and put in memory in create template cache
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	return db, nil
}

// TODO -
/*
	- passing context everywhere
	- logging level correction
	- admin ui bug solve
	- writing tests

*/
