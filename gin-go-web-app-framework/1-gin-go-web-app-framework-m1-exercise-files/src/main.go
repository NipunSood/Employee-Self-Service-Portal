package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
)

const port = "3000"

func main() {
	var templates *template.Template
	var tmplCh = make(chan *template.Template)

	go func(ch <-chan *template.Template) {
		for tmpl := range tmplCh {
			templates = tmpl
		}
	}(tmplCh)

	go watchTemplateFolder(tmplCh)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		if len(path) == 0 {
			path = "index"
		}
		templates.Lookup(path+".html").Execute(w, nil)
	})
	http.Handle("/public/", http.FileServer(http.Dir(".")))
	go http.ListenAndServe(":"+port, nil)

	fmt.Println("Server is running on port " + port)
	var a string
	fmt.Scanln(&a)
}

func parseTemplates(ch chan<- *template.Template) {
	result := template.New("")

	const basePath = "templates/"

	template.Must(result.ParseGlob(basePath + "includes/**"))
	template.Must(result.ParseGlob(basePath + "views/**"))

	ch <- result
}
func watchTemplateFolder(ch chan<- *template.Template) {
	updateTimes := map[string]time.Time{}
	const basePath = "templates"
	var shouldUpdate bool
	var scanFolder func(string)

	scanFolder = func(folder string) {
		d, _ := os.Open(folder)
		defer d.Close()
		fis, _ := d.Readdir(-1)
		for _, fi := range fis {
			if fi.IsDir() && !shouldUpdate {
				scanFolder(folder + "/" + fi.Name())
			} else {
				if updateTimes[folder+"/"+fi.Name()].Unix() < fi.ModTime().Unix() {
					shouldUpdate = true
					updateTimes[folder+"/"+fi.Name()] = fi.ModTime()
				}
			}
		}
	}

	for {
		time.Sleep(500 * time.Millisecond)
		shouldUpdate = false

		scanFolder(basePath)

		if shouldUpdate {
			parseTemplates(ch)
		}
	}
}
