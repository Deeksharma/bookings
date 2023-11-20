package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Deeksharma/bookings/internal/config"
	"github.com/Deeksharma/bookings/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var functions = template.FuncMap{
	"humanDate": HumanDate,
	"getBool":   GetBool,
} // specifies functions that will be available to our golang templates

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer set the config for template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate return time in dd-mm-yyyy format
func HumanDate(t time.Time) string {
	return t.Format("02-01-2006")
}

// GetBool returns string boolean for an integer
func GetBool(val int) string {
	if val == 0 {
		return "false"
	}
	return "true"
}

func AddDefaultData(data *models.TemplateData, r *http.Request) *models.TemplateData {
	data.Flash = app.SessionManager.PopString(r.Context(), "flash") // this will put something in the session until the next time a page is displayed and then its removed automatically
	data.Warning = app.SessionManager.PopString(r.Context(), "warning")
	data.Error = app.SessionManager.PopString(r.Context(), "error")
	data.CSRFToken = nosurf.Token(r)
	if app.SessionManager.Exists(r.Context(), "user_id") {
		data.IsAuthenticated = 1
	}
	return data
}

func Template(writer http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) error {
	var templateCache map[string]*template.Template
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}
	// get the template cache from app config
	//templateCache, err := CreateTemplateCache()
	//if err != nil {
	//	log.Fatal(err)
	//}
	temp, ok := templateCache[tmpl]
	if !ok {
		log.Println("could not get template from template cache")
		return errors.New("cant get templates from cache")
	}

	buf := new(bytes.Buffer)
	data = AddDefaultData(data, r)
	_ = temp.Execute(buf, data) // agar koi data hota h toh woh yaha store hota h

	_, err := buf.WriteTo(writer)
	if err != nil {
		fmt.Println("Error writing template to browser: ", err)
		return err
	}
	return nil

	//parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	//err = parsedTemplate.Execute(writer, nil)
	//if err != nil {
	//	fmt.Println("Error rendering template")
	//	return
	//}
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) { // find the template and the layouts associated with it
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
