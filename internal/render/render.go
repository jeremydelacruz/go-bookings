package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jeremydelacruz/go-bookings/internal/config"
	"github.com/jeremydelacruz/go-bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

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

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) {
	var templateCache map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		var err error
		templateCache, err = CreateTemplateCache()
		if err != nil {
			log.Fatal("could not create template cache")
		}
	}

	// get requested template from cache
	t, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}
	log.Println("rendered from template cache")

	// use a buffer here just as another potential point of error handling
	buf := new(bytes.Buffer)
	data = addDefaultData(data, r)
	err := t.Execute(buf, data)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	log.Println("creating template cache")
	templateCache := map[string]*template.Template{}

	// filepath.glob returns the full path of all template files
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return templateCache, err
	}

	// store if there are existing layout files
	matches, err := filepath.Glob("./templates/*.layout.tmpl")
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
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = templateSet
	}

	return templateCache, nil
}
