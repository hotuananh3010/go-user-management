package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func main() {
	router := httprouter.New()
	router.GET("/", Home)
	router.GET("/home", Home)
	router.GET("/users", users)
	router.GET("/users/add", userAdd)
	router.POST("/users/add", userAdd)
	router.GET("/users/edit/:id", userEdit)
	router.POST("/users/edit/:id", userEdit)
	log.Fatal(http.ListenAndServe(":8080", router))

	// router := gin.Default()
	// router.GET("/albums", getAlbums)

	// router.Run("localhost:8080")
}

type Page struct {
	Title string
	Body  []byte
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

type UserModel struct {
	id      int
	code    string `validate:"required"`
	name    string `validate:"required"`
	program string `validate:"required"`
}

func users(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data []UserModel

	db, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer db.Close()

	row, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var code string
		var name string
		var program string
		row.Scan(&id, &code, &name, &program)
		item := UserModel{
			id:      id,
			code:    code,
			name:    name,
			program: program,
		}
		data = append(data, item)
		//log.Println("Student: ", code, " ", name, " ", program)
	}
	renderTemplate(w, "templates/users/list.html", data)
}

func userAdd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var data []UserModel
	db, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer db.Close()
	if r.Method == "GET" {
		row, err := db.Query("SELECT * FROM users")
		if err != nil {
			log.Fatal(err)
		}
		defer row.Close()
		for row.Next() { // Iterate and fetch the records from result cursor
			var id int
			var code string
			var name string
			var program string
			row.Scan(&id, &code, &name, &program)
			item := UserModel{
				id:      id,
				code:    code,
				name:    name,
				program: program,
			}
			data = append(data, item)
			//log.Println("Student: ", code, " ", name, " ", program)
		}
		renderTemplate(w, "templates/users/add.html", data)
	} else if r.Method == "POST" {
		validate = validator.New(validator.WithRequiredStructEnabled())

		user := &UserModel{
			code:    r.FormValue("code"),
			name:    r.FormValue("name"),
			program: r.FormValue("program"),
		}
		err := validate.Struct(user)
		if err != nil {
			log.Fatal(err)
			http.Redirect(w, r, "/users/add", http.StatusFound)
		}

		//fmt.Fprintf(w, "Hi there, I love %s!", r.Method)
		log.Println("Inserting student record ...")
		insertStudentSQL := `INSERT INTO users(code, name, program) VALUES (?, ?, ?)`
		statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = statement.Exec(user.code, user.name, user.program)
		if err != nil {
			log.Fatalln(err.Error())
		}
		http.Redirect(w, r, "/users", http.StatusFound)
		return
	}
}

func userEdit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var data UserModel
	id := ps.ByName("id")
	if _, err := strconv.Atoi(id); err == nil {
		log.Printf("%q looks like a number.\n", id)
	}

	db, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer db.Close()

	if r.Method == "GET" {
		row, err := db.Query("SELECT * FROM users WHERE id = ?", id)
		if err != nil {
			log.Fatal(err)
		}
		defer row.Close()
		for row.Next() { // Iterate and fetch the records from result cursor
			var id int
			var code string
			var name string
			var program string
			row.Scan(&id, &code, &name, &program)
			data = UserModel{
				id:      id,
				code:    code,
				name:    name,
				program: program,
			}
		}
		renderTemplate(w, "templates/users/edit.html", data)
	} else if r.Method == "POST" {
		id := r.FormValue("id")
		code := r.FormValue("code")
		name := r.FormValue("name")
		program := r.FormValue("program")
		//fmt.Fprintf(w, "Hi there, I love %s!", r.Method)
		log.Println("Inserting student record ...")
		insertStudentSQL := `UPDATE users SET code = ? , name = ?, program = ? WHERE id = ?`
		statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = statement.Exec(code, name, program, id)
		if err != nil {
			log.Fatalln(err.Error())
		}
		http.Redirect(w, r, "/users", http.StatusFound)
		return
	}

}
