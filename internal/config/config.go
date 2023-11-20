package config

import (
	"github.com/Deeksharma/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

// config should not import each other for avoiding import cycle

// AppConfig holds the application config
type AppConfig struct {
	UseCache       bool
	TemplateCache  map[string]*template.Template
	InfoLog        *log.Logger // now every part of application will have same logger
	ErrorLog       *log.Logger
	InProduction   bool
	SessionManager *scs.SessionManager // currently it uses cookies to store session, we can use redis as well later
	MailChannel chan models.MailData // mail channel
}
