package render

import (
	"encoding/gob"
	"github.com/Deeksharma/bookings/internal/config"
	"github.com/Deeksharma/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var sessionManager *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	// change it to true when in production
	testApp.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true // cookie will persist even after the browser is closed
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false // since localhost is not secure

	testApp.SessionManager = sessionManager

	app = &testApp

	os.Exit(m.Run())
}

type myHttpWriter struct {}

func (w *myHttpWriter) Header() http.Header{
	var header http.Header
	return header
}
func (w *myHttpWriter) Write(b []byte) (int, error){
	return len(b), nil
}
func (w *myHttpWriter) WriteHeader(statusCode int) {
}
