package controller

import (
	"encoding/json"
	"example/server/database"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var user model.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		fmt.Println("error occured", err)
	}
	fmt.Println(user)
	err = myvalidator(user) // Validating fields
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("i am validated")
	hp, _ := HashPassword(user.Password) //  Hashing password
	user.Password = hp
	user = database.CreateUser(&user)
	ruser := model.Mapping(user)
	ruser.Age = age.Age(*user.DateOfBirth)
	json.NewEncoder(w).Encode(ruser)
}

//Updating user
func UpdateUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var user model.User
	ID := mux.Vars(req)
	x := ID["id"]
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
	perpage := 2
	var total int64
	var x int
	DB.Raw(sql).Count(&total).Scan(x)
	sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, perpage, (page-1)*perpage)
	DB.Raw(sql).Scan(&ruser)

	if len(ruser) == 0 {
		w.WriteHeader(404)
		fmt.Fprintln(w, "NO DATA FOUND")
		return
	}
	pages := model.Paginate{
		Data:      ruser,
		Total:     total,
		Page:      page,
		Last_page: math.Ceil(float64(total / int64(perpage))),
	}
	json.NewEncoder(w).Encode(pages)
}

func Login(w http.ResponseWriter, res *http.Request) {
	fmt.Printf("i am login")
	var cred model.Cred
	var user model.User
	DB := database.GetDb()

	err := json.NewDecoder(res.Body).Decode(&cred)
	fmt.Println(cred.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	DB.First(&user, "email= ?", cred.Email)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
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
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	fmt.Fprintln(w, "you are logged in")
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
