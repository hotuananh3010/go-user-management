package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/go-playground/validator/v10"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate
var db *sql.DB

func main() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	router := httprouter.New()
	db, _ = sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer db.Close()

	router.GET("/", Home)
	router.GET("/home", Home)
	router.GET("/users", users)
	router.GET("/users/add", UserAdd)
	router.POST("/users/add", UserStore)
	router.GET("/users/edit/:id", UserEdit)
	router.POST("/users/edit/:id", UserUpdate)
	router.POST("/users/delete", UserDelete)
	log.Fatal(http.ListenAndServe(":8080", router))
}

//var templates = template.Must(template.ParseFiles("templates/home.html", "templates/about.html", "templates/users/list.html", "templates/users/edit.html", "templates/users/add.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p any) {
	t, err := template.New(path.Base(tmpl)).ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = t.Execute(w, p)
	//err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//log.Println("Inserting student record ...")
	data := struct {
		Title string
		Body  string
	}{
		"Home",
		"HOME body",
	}
	renderTemplate(w, "templates/home.html", data)
}
