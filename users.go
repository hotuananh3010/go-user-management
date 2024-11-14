package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

func users(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data []UserModel
	// db, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	// defer db.Close()

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
			Id:      id,
			Code:    code,
			Name:    name,
			Program: program,
		}
		data = append(data, item)
	}

	var pageModel = PageModel{
		Title: "Users",
		Data:  data,
	}
	renderTemplate(w, "templates/users/list.html", pageModel)
}

func UserAdd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var data []UserModel

	flashMessages, _ := GetFlash(w, r)

	var pageModel = PageModel{
		Title:         "Add User",
		FlashMessages: flashMessages,
		Data:          data,
	}
	renderTemplate(w, "templates/users/add.html", pageModel)
}

func UserStore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := &UserModel{
		Code:    r.FormValue("code"),
		Name:    r.FormValue("name"),
		Program: r.FormValue("program"),
	}
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			SetFlash(w, r, err.Field()+" "+err.Tag())
		}
		http.Redirect(w, r, "/users/add", http.StatusFound)
	} else {
		insertStudentSQL := `INSERT INTO users(code, name, program) VALUES (?, ?, ?)`
		statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = statement.Exec(user.Code, user.Name, user.Program)
		if err != nil {
			log.Fatalln(err.Error())
		}
		http.Redirect(w, r, "/users", http.StatusFound)
	}
}

func UserEdit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var data UserModel
	id := ps.ByName("id")
	if _, err := strconv.Atoi(id); err == nil {
		log.Printf("%q looks like a number.\n", id)
	}

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
			Id:      id,
			Code:    code,
			Name:    name,
			Program: program,
		}
	}
	flashMessages, _ := GetFlash(w, r)

	var pageModel = PageModel{
		Title:         "Edit User",
		FlashMessages: flashMessages,
		Data:          data,
	}

	renderTemplate(w, "templates/users/edit.html", pageModel)
}

func UserUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if _, err := strconv.Atoi(id); err == nil {
		log.Printf("%q looks like a number.\n", id)
	}

	user := &UserModel{
		Code:    r.FormValue("code"),
		Name:    r.FormValue("name"),
		Program: r.FormValue("program"),
	}
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			SetFlash(w, r, err.Field()+" "+err.Tag())
		}
		http.Redirect(w, r, "/users/edit/"+id, http.StatusFound)
	} else {
		insertStudentSQL := `UPDATE users SET code = ? , name = ?, program = ? WHERE id = ?`
		statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = statement.Exec(user.Code, user.Name, user.Program, id)
		if err != nil {
			log.Fatalln(err.Error())
		}
		http.Redirect(w, r, "/users", http.StatusFound)
	}

}

func UserDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := r.FormValue("id")
	log.Printf("%q looks like a number.\n", id)
	if _, err := strconv.Atoi(id); err != nil {
		SetFlash(w, r, "User not found")
		http.Redirect(w, r, "/users", http.StatusFound)
	}

	insertStudentSQL := `DELETE FROM users WHERE id = ?`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// // This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(id)
	if err != nil {
		log.Fatalln(err.Error())
	}
	http.Redirect(w, r, "/users", http.StatusFound)

}
