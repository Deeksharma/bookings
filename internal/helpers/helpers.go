package helpers

import (
	"fmt"
	"github.com/Deeksharma/bookings/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

// We can have client errors and server errors

func ClientError(w http.ResponseWriter, status int) {
	app.ErrorLog.Println("Client error with status of ", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) // later we'll have a log file or send some email
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(request *http.Request) bool {
	exists := app.SessionManager.Exists(request.Context(), "user_id")
	return exists
}
