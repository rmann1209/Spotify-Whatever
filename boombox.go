package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/fs"
	"log"
	"net/http"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique"` //specifies as unique, blocks repeats
	Password  string
	FirstName string
	LastName  string
}

type SongReview struct {
	gorm.Model
	SongTitle string // song Title for song Review
	Rating    int    // Rating out of 5
	Comment   string // User comment on song
	Username  string // for the user
}

type AlbumReview struct {
	gorm.Model
	AlbumTitle string // album Title for song Review
	Rating     int    // Rating out of 5
	Comment    string // User comment on song
	Username   string // for the user
}

var db *gorm.DB
var static embed.FS

var activeUsername string = ""

func main() {
	//db -> database
	var err error
	db, err = gorm.Open(sqlite.Open("usersbase.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to open database.")
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to automigrate")
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/home", HomeHandler)
	r.HandleFunc("/signup", SignUpHandler)
	r.HandleFunc("/login", LoginHandler)
	http.Handle("/", r)
	webapp, err := fs.Sub(static, "static")
	if err != nil {
		fmt.Println(err)
	}
	r.PathPrefix("/").Handler(http.FileServer(http.FS(webapp)))
	//starts logger
	r.Use(logger)

	fmt.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil) //DO NOT CHANGE THIS
	if err != nil {
		log.Fatalln(err)
	}
}
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(*r)
	if r.Method == "GET" {
		fmt.Println("GET REQUEST RECEIVED ON SIGNUP")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if r.Method != "POST" {
		fmt.Printf("POST REQUEST NOT RECEIVED ON SIGNUP %s received instead", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	enableCors(&w) //so that the FE can access
	var signupUser User
	length := r.ContentLength
	if length > 0 {
		json.NewDecoder(r.Body).Decode(&signupUser)
	}
	username := signupUser.Username
	password := signupUser.Password
	firstName := signupUser.FirstName
	lastName := signupUser.LastName

	result := db.Create(&User{
		Username:  username,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	})

	if result.Error != nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	} else {

		fmt.Printf("New username: %s  and password %s", username, password)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	enableCors(&w) //so that the FE can access
	var loginUser User
	length := r.ContentLength
	if length > 0 {
		json.NewDecoder(r.Body).Decode(&loginUser)
	}
	username := loginUser.Username
	password := loginUser.Password

	user := User{}

	result := db.Where("username = ? AND password = ?", username, password).First(&user)
	if result.Error != nil {
		http.Error(w, "Username or password is incorrect", http.StatusBadRequest)
		return
	} else {
		//make logged in user the activeUser
		activeUsername = loginUser.Username
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	/* if r.Method == "GET" {
		tmpl, err := template.ParseFiles("homemock.html") //direct to the file
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}
	*/
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	result := r.FormValue("action")
	if result == "Sign up" {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
	} else if result == "Login" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)

	}

}

// begin logger - DND
func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func enableCors(w *http.ResponseWriter) { //this function enables Cors which may be used to link FE and BE
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}
