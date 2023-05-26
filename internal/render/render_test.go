package render

import (
	"net/http"
	"testing"

	"github.com/jeremydelacruz/go-bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getRequestWithSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")
	res := addDefaultData(&td, r)
	if res.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc
	r, err := getRequestWithSession()
	if err != nil {
		t.Error(err)
	}

	var tw testWriter
	err = Template(&tw, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser", err)
	}

	err = Template(&tw, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getRequestWithSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/test-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
