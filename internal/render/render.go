package render

import (
	"bookings/internal/config"
	"bookings/internal/models"
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"

	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, templateData *models.TemplateData) {
	var templateCache map[string]*template.Template
	if app.UseCache {
		// get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	t, ok := templateCache[tmpl] // ok will be true if the tmpl exists
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData, r)

	_ = t.Execute(buf, templateData)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = templateSet
	}
	return myCache, nil
}
