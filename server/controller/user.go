package controller

import (
	"encoding/json"
	"example/server/database"
	"example/server/middleware"

	"example/server/model"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	age "github.com/bearbin/go-age"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Validation
func myvalidator(x model.User) error {
	validate := validator.New()
	return validate.Struct(x)
}

//Hashing Password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Signup API
func Signup(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var user model.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		fmt.Println("error occured", err)
	}
	// Validating fields
	err = myvalidator(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hp, _ := HashPassword(user.Password) //  Hashing password
	user.Password = hp
	user = database.CreateUser(&user)
	if user.ID == 0 {
		http.Error(w, "UserAlreadyFound", http.StatusBadRequest)
		return
	}
	ruser := model.Mapping(user)
	err = json.NewEncoder(w).Encode(ruser)
	if err != nil {
		fmt.Println(err.Error())
	}
}

//Updating user
func UpdateUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var user model.User
	id := mux.Vars(req)
	x := id["id"]
	user = database.GetUserById(x)
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		fmt.Println("error occured", err)
	}
	err = myvalidator(user) // Validating fields
	if err != nil {
		x := err.Error()
		if x == "Key: 'User.Password' Error:Field validation for 'Password' failed on the 'max' tag" {
		} else {
			http.Error(w, x, http.StatusBadRequest)
			return
		}
	}
	hp, _ := HashPassword(user.Password) //  Hashing password
	user.Password = hp
	user = database.UpdateUser(&user)
	ruser := model.Mapping(user)
	json.NewEncoder(w).Encode(ruser)
}

//Delete User

func DeleteUser(w http.ResponseWriter, req *http.Request) {
	ID := mux.Vars(req)
	x := ID["id"]
	err := database.DeleteUser(x)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintln(w, "NO USER TO BE DELETED")
		return
	}
	fmt.Fprintln(w, "USER DELETED")
}

//read
func Read(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	DB := database.GetDb()
	var ruser []model.ResUser
	//Parameters
	archived := req.URL.Query().Get("archived")
	id := req.URL.Query().Get("id")
	firstname := req.URL.Query().Get("firstname")
	email := req.URL.Query().Get("email")
	sort := req.URL.Query().Get("sort")
	order := req.URL.Query().Get("order")

	sql := "SELECT * FROM users where archived=0"
	if archived == "true" {
		sql = "SELECT * FROM users where archived=1"
	}
	if id != "" {
		sql = fmt.Sprintf("%s and id = %s", sql, id)
	}
	if firstname != "" {
		sql = fmt.Sprintf("%s and first_name LIKE  '%%%s%%'", sql, firstname)
	}
	if email != "" {
		sql = fmt.Sprintf("%s and email LIKE  '%%%s%%'", sql, email)
	}
	if sort != "" {
		if order == "" {
			order = "asc"
		} else {
			sql = fmt.Sprintf("%s order by %s %s", sql, sort, order)
		}
	}
	// pagination
	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	perpage := 10
	var total int
	sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, perpage, (page-1)*perpage)
	DB.Raw(sql).Scan(&ruser)
	total = len(ruser)
	if len(ruser) == 0 {
		w.WriteHeader(404)
		fmt.Fprintln(w, "NO DATA FOUND")
		return
	}
	for i, _ := range ruser {
		dateString := ruser[i].DateOfBirth
		date, _ := time.Parse("2006-01-02", dateString)
		ruser[i].Age = age.Age(date)
	}
	pages := model.Paginate{
		Data:      ruser,
		Total:     total,
		Page:      page,
		Last_page: math.Ceil(float64(total / (perpage))),
	}
	json.NewEncoder(w).Encode(pages)
}

// Login
func Login(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var cred model.Cred
	var user model.User
	DB := database.GetDb()
	type message struct {
		Info   string
		Status int
	}
	var mess message
	x, _ := req.Cookie("token")
	if x != nil {
		mess = message{
			Info:   "Already logged in",
			Status: 401,
		}
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(mess)

		return

	}
	err := json.NewDecoder(req.Body).Decode(&cred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	DB.First(&user, "email= ?", cred.Email)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		mess = message{
			Info:   "invalid email or password",
			Status: 401,
		}
		json.NewEncoder(w).Encode(mess)
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &model.Claims{
		Email: cred.Email,
		StandardClaims: jwt.StandardClaims{

			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	TokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("token", TokenString)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    TokenString,
		Expires:  expirationTime,
		HttpOnly: x.Secure,
	})
	fmt.Println(req.Cookie("token"))
	mess = message{
		Info:   "you are logged in",
		Status: 200,
	}
	json.NewEncoder(w).Encode(mess)
}

//Logout
func Logout(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("token")
	if err != nil {
		http.Error(w, "Already logged out", http.StatusBadRequest)
		return
	}

	cookie = &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	json.NewEncoder(w).Encode("logged out")
}
