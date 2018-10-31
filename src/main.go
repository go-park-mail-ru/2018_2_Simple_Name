package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const letterBytes string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type User struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Nick     string `json:"nick"`
	Email    string `json:"email"`
	Password string `json:"password"`
	KeyWord  string `json:"-"`
	Score    int    `json:"score"`
	Age      int    `json:"age"`
}

type Users map[string]User

var users = make(Users)

type Sessions map[string]string

var sessions = make(Sessions)

func CORSsettings(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, User-Agent, Cache-Control, Accept, X-Requested-With, If-Modified-Since")
		if r.Method == http.MethodOptions{
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func main() {

	http.HandleFunc("/signin", CORSsettings(signinHandler))
	http.HandleFunc("/signup", CORSsettings(signupHandler))
	http.HandleFunc("/profile", CORSsettings(profileHandler))
	http.HandleFunc("/leaders", CORSsettings(leadersHandler))
	http.HandleFunc("/islogged", CORSsettings(islogged))

	fmt.Println("starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("cannot listen: %s", err)
	}
}

func leadersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC(), "Request from", r.URL.String())
	fmt.Println("Method", r.Method)

	Leaders := map[int]User{
		0: User{
			Nick:  "GRe12",
			Score: 4321,
			Age:   12,
		},
		1: User{
			Nick:  "wasaW2",
			Score: 43121,
			Age:   13,
		},
		2: User{
			Nick:  "Feesfs",
			Score: 432441,
			Age:   77,
		},
	}

	w.Header().Set("Content-Type", "application/json")

	resp, _ := json.Marshal(Leaders)
	w.Header().Set("Status-Code", "200")

	w.Write(resp)

	return
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC(), "Request from", r.URL.String())
	fmt.Println("Method", r.Method)

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//user, err := getFormReq(r)

	user, err := getJSONReq(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.Email != "" && user.Password != "" {
		if !checkExist(*user) {

			user.KeyWord = RandStringBytesRmndr()
			id, _ := addUser(*user)
			session := new(http.Cookie)
			session.Name = "session_id"
			session.Value = uidGen()
			session.Expires = time.Now().Add(time.Second)
			session.HttpOnly = true
			http.SetCookie(w,session)
			sessions[session.Value] = id
			users[user.Email] = *user

			w.WriteHeader(http.StatusCreated)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC(), "Request from", r.URL.String())
	fmt.Println("Method", r.Method)

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//user, err := getFormReq(r)

	user, err := getJSONReq(r)

	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(user.Email)
	fmt.Println(user.Password)

	if validateUser(*user) {
		session := new(http.Cookie)
		session.Name = "session_id"
		session.Value = uidGen()
		session.Expires = time.Now().Add(time.Hour)
		session.HttpOnly = true
		http.SetCookie(w, session)

		sessions[session.Value] = user.Email

		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC(), "Request from", r.URL.String())
	fmt.Println("Method", r.Method)

	//Online, id := loggedIn(r)
	var id = ""

	if r.Method == http.MethodGet {

		//if !Online {
		//
		//	return
		//}
		userJson, err := json.Marshal(users[id])
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Status-Code", "200")

		w.Write(userJson)
		return
	}

	if err := uploadFileReq(id, r); err != nil {

		return
	}

	user := users[id]
	data, err := getFormReq(r)
	if err != nil {
		return
	}
	if data.Nick != "" {
		user.Nick = data.Nick
	}
	if data.Password != "" {
		user.Password = data.Password
	}
	users[id] = user

	return

}

func islogged(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC(), "Request from", r.URL.String())
	fmt.Println("Method", r.Method)

	val := r.Cookies()

	for i := 0; i < len(val); i++{
		_, ok := sessions[val[i].Value]
		if !ok {
			continue
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}

	}
	w.WriteHeader(http.StatusBadRequest)
	return
}


func addUser(user User) (string, bool) {
	users[user.Email] = user
	fmt.Println(time.Now().UTC(), "Added user", user)

	return user.Email, true
}

func checkExist(user User) bool {
	if _, ok := users[user.Email]; ok {
		return true
	}
	return false
}

func validateUser(user User) bool {
	fmt.Println(users)


	fmt.Println("User exist: ", users[user.Email].Email)
	fmt.Println("User: ", users[user.Email].Password)

	fmt.Println("User: ", user.Email)
	fmt.Println("Password: ", user.Password)
	if mapUser, ok := users[user.Email]; ok {
		if user.Password == mapUser.Password {
			return true
		}
	}
	return false
}

func RandStringBytesRmndr() string {
	n := 10
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func getJSONReq(r *http.Request) (*User, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		fmt.Println("Ошибка чтения 1: ", err.Error())
		return nil, err
	}

	user := new(User)

	err = json.Unmarshal(body, user)

	if err != nil {
		fmt.Println("Ошибка чтения 2: ", err.Error())
		return nil, err
	}

	return user, nil
}

func getFormReq(r *http.Request) (*User, error) {
	user := new(User)
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")
	user.Name = r.FormValue("name")
	user.LastName = r.FormValue("last_name")
	user.Nick = r.FormValue("nick")

	return user, nil
}

func uploadFileReq(fileName string, r *http.Request) error {
	if err := r.ParseMultipartForm(32 << 20); nil != err {
		fmt.Println(err.Error())
		return err
	}

	file, _, err := r.FormFile("my_file")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer file.Close()

	dst, err1 := os.Create(filepath.Join("../src/static/media", fileName))

	if err1 != nil {
		fmt.Println(err1.Error())
		return err
	}

	io.Copy(dst, file)
	return nil
}

func uidGen() string {
	n := 15
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
	}
	return string(b)
}
