package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/jeremydelacruz/go-bookings/internal/config"
	"github.com/jeremydelacruz/go-bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"
var functions = template.FuncMap{}

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// addDefaultData adds data that should be present on every page
func addDefaultData(data *models.TemplateData, r *http.Request) *models.TemplateData {
	data.CSRFToken = nosurf.Token(r)
	data.Flash = app.Session.PopString(r.Context(), "flash")
	data.Warning = app.Session.PopString(r.Context(), "warning")
	data.Error = app.Session.PopString(r.Context(), "error")
	return data
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) error {
	var templateCache map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		var err error
		templateCache, err = CreateTemplateCache()
		if err != nil {
			return fmt.Errorf("run: failed creating template cache: %w", err)
		}
	}

	// get requested template from cache
	t, ok := templateCache[tmpl]
	if !ok {
		return fmt.Errorf("RenderTemplate: failed fetching from template cache")
	}

	// use a buffer here just as another potential point of error handling
	buf := new(bytes.Buffer)
	data = addDefaultData(data, r)
	err := t.Execute(buf, data)
	if err != nil {
		return fmt.Errorf("RenderTemplate: failed executing parsed template: %w", err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		return fmt.Errorf("RenderTemplate: failed writing buffer to response writer: %w", err)
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	log.Println("creating template cache")
	templateCache := map[string]*template.Template{}

	// filepath.glob returns the full path of all template files
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return templateCache, err
	}

	// store if there are existing layout files
	matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
	if err != nil {
		return templateCache, err
	}

	// iterate through each page, also parsing all layouts with each page
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = templateSet
	}

	return templateCache, nil
}
