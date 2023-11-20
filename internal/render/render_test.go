package render

import (
	"github.com/Deeksharma/bookings/internal/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var data models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	sessionManager.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&data, r)
	if result.Flash != "123" {
		t.Error("flash value 123 not found")
	}
}

func TestRenderTemplate(t *testing.T)  {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myHttpWriter
	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to browser",err)
	}

	err = Template(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("non existent template got written",err)
	}
}

func TestNewTemplates(t *testing.T){
	NewRenderer(app)
}

func TestCerateTemplateCache(t *testing.T){
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error){
	r, err := http.NewRequest("GET", "/some-urls", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = sessionManager.Load(ctx, r.Header.Get("X-Session")) // session manager should have context

	r = r.WithContext(ctx)
	return r, nil
}
