package controller

import (
	"encoding/json"
	//"fmt"
	Database "github.com/nekonotes/database"
	helper "github.com/nekonotes/helper"
	"github.com/nekonotes/middleware"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

var db, notes_db = Database.ConnecttoDB()
func Homepage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicatioan/json")
	json.NewEncoder(w).Encode("hello world")
}
func verifypassword(providedPassword string, userPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	if err != nil {
		check = false
	}
	return check
}
func Register(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is OPTIONS (preflight request)
	if r.Method == "OPTIONS" {
		// Respond with a success status for preflight requests
		w.WriteHeader(http.StatusOK)
		return
	}

	// Decode JSON data from request body
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	// Retrieve values from decoded JSON
	name := data["name"]
	email := data["email"]
	password := data["password"]
	//check if the any user exist by that user name or is between size 3 and 100
	if middleware.Name_check(name, db) != "done" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middleware.Name_check(name, db))
		return
	}

	//hash the password
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	//check for valid email format
	if middleware.EmailCheck(email, db) != "done" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middleware.EmailCheck(email, db))
		return
	}

	token, refresh_token, _ := helper.GenerateAlltoken(email, name, "user", db)
	var v Database.User = Database.User{Name: name, Email: email, Password: string(hashedpassword), User_type: "user", Token: token, Refresh_token: refresh_token}
	db.Create(&v)
	helper.SetCookieHandler(w, r, token, refresh_token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("created")
}
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		// Respond with a success status for preflight requests
		w.WriteHeader(http.StatusOK)
		return
	}

	// Decode JSON data from request body
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	email := data["email"]
	password := data["password"]

	var users []Database.User
	db.Find(&users, "email = ?", email)
	for _, user := range users {
		if user.Email == email {
			if verifypassword(user.Password, string(password)) {
				// email and password matched
				token, refresh_token, _ := helper.GenerateAlltoken(email, user.Name, "user", db)
				helper.GenerateAlltoken(token, refresh_token, email, db)
				helper.SetCookieHandler(w, r, token, refresh_token)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("verified")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("wrong password try again")
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("invalid user name try again")
}
func CreateNotes(w http.ResponseWriter, r *http.Request) {
	// Decode JSON data from request body
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	// read the cookie
	value, err := helper.ReadCookie(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	//read the jwt token
	claims, err := helper.IsAuthorized(value["token"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	notes := data["notes"]
	title := data["title"]
	name := claims.Name
	var v Database.User_Notes = Database.User_Notes{Name: name, Notes: notes, Title: title}
	notes_db.Create(&v)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("created")
}
func GetNotesTitle(w http.ResponseWriter, r *http.Request) {
	
	// read the cookie
	value, err := helper.ReadCookie(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	//read the jwt token
	claims, err := helper.IsAuthorized(value["token"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	name := claims.Name
	var response []map[string]uint
	var notes []Database.User_Notes
	notes_db.Find(&notes, "name = ?", name)
	for _, user := range notes {
		if user.Name == name {
			temp := make(map[string]uint)
			temp[user.Title]=user.ID
			response = append(response, temp)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func GetNote(w http.ResponseWriter, r *http.Request){
	// Decode JSON data from request body
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	// read the cookie
	value, err := helper.ReadCookie(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	//read the jwt token
	claims, err := helper.IsAuthorized(value["token"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	name := claims.Name
	id,_:= strconv.Atoi( data["id"])
	var notes []Database.User_Notes
	notes_db.Find(&notes, "Id = ?", id)
	for _, user := range notes {
		if user.ID ==  uint(id) && user.Name == name {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}
}
func DeleteUser(w http.ResponseWriter, r *http.Request){
	// Decode JSON data from request body
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	// read the cookie
	value, err := helper.ReadCookie(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	//read the jwt token
	_, err = helper.IsAuthorized(value["token"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	var user Database.User
	email := data["email"]

    if err := db.Where("email = ?", email).First(&user).Error; err != nil {
        w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("invalid email try again")
		return 
    }
	if !verifypassword(user.Password, string(data["password"])){
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("incorrect password")
		return
	}
	
    // Delete the user from the database
    if err := db.Delete(&user).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("error occured while deleting")
		return 
    }
	name := user.Name
	var notes []Database.User_Notes
	notes_db.Find(&notes, "name = ?", name)
	for _, user := range notes {
		if user.Name == name {
			notes_db.Delete(&user)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("deleted")
}
func DeleteNote(w http.ResponseWriter, r *http.Request){
	// Decode JSON data from request body
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {

		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	// read the cookie
	value, err := helper.ReadCookie(r)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	//read the jwt token
	claims, err := helper.IsAuthorized(value["token"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	name := claims.Name
	id,_:= strconv.Atoi( data["id"])
	var notes []Database.User_Notes
	notes_db.Find(&notes, "Id = ?", id)
	for _, user := range notes {
		if user.ID ==  uint(id) && user.Name == name {

			notes_db.Delete(&user)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("deleted")
			return 
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("not found")
}