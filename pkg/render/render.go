package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/karthikkalarikal/golangLogin/pkg/config"
	"github.com/karthikkalarikal/golangLogin/pkg/models"
)

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

var app *config.Appconfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.Appconfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get the template cache from the app config

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	//render the template

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}
func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	//get all of the files named *.page.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}
	//range through all files ending with *.page.html
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
		}
		if err != nil {
			return myCache, err
		}
		myCache[name] = ts
	}
	return myCache, nil

}

// parsedTemplate, err := template.ParseFiles("./templates/"+tmpl, "templates/base.layout.html")
// if err != nil {
// 	http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
// 	return
// }
// td = AddDefaultData(td)
// err = parsedTemplate.Execute(w, td)
// if err != nil {
// 	http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
// 	return
// }
// // var tc = make(map[string]*template.Template)

// func RenderTemplateTest(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error
// 	//check to see if we already have the template
// 	_, inMap := tc[t]
// 	if !inMap {
// 		//need to create the template
// 		err = CreateTemplateCache(t)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		// we have the template
// 		log.Println("using cached template")
// 	}

// 	tmpl = tc[t]
// 	err = tmpl.Execute(w, nil)
// }
// func CreateTemplateCache(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.layout.html",
// 	}
// 	//parse the template
// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}
// 	tc[t] = tmpl
// 	return nil
// }
