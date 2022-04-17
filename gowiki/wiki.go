package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

const (
	sep string = string(os.PathSeparator)

	// Please note that folders must be within the working environment (learningo/gowiki/ in this case)
	tmplFolder string = sep + "pages" + sep
	htmlFolder string = sep + "html" + sep
)

var (
	// env(): Function to get the current environment (ie, .../gowiki/)
	// Note: This could be a simple variable as well, avoiding multiple execution to get the same output.
	env = func() string {
		path, err := os.Getwd()
		panicError(err)
		return path
	}

	// panicError Simple error handler that will panic if err is not null and print its content
	panicError = func(err error) {
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}
	}
	viewPath  = env() + htmlFolder + "view.html"
	editPath  = env() + htmlFolder + "edit.html"
	templates = template.Must(template.ParseFiles(viewPath, editPath))
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := env() + tmplFolder + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

// Load a page from a text file (.txt)
func loadPage(title string) (*Page, error) {
	filename := env() + tmplFolder + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// viewHandler will allow users to view a wiki page. It will handle URLs prefixed with "/view/".
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	render("view")(w, p)
}

// The function editHandler loads the page (or, if it doesn't exist, create an empty Page struct), and displays an HTML form.
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	render("edit")(w, p)
	/*
		Although previous execution is less verbose, one could call the same function as:
		renderEdit = render("edit")

		And then call it like so:
		renderEdit(w,p)

		This approach obviously increments verbosity but makes the code more clear and granular.
	*/
}

// Don't ask me why, but I wanted to implement some currying in Go, hence first define the template type and then execute it.
func render(templ string) func(w http.ResponseWriter, p *Page) {
	templ = templ + ".html"
	return func(w http.ResponseWriter, p *Page) {
		err := templates.ExecuteTemplate(w, templ, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// saveHandler will handle the submission of forms located on the edit pages
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// Function that uses the validPath expression to validate path and extract the page title.
// Note: Not used function, but leaving it here to showcase the get title implementation.
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
