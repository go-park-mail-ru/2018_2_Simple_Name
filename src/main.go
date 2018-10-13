package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	Email    string
	Password string `json: "-"`
	KeyWord  string `json: "-"`
	Score    int    `json: "score"`
	Age      int    `json: "age"`
}

type Users map[string]User

var users = make(Users)

type Sessions map[string]string

var sessions = make(Sessions)

func CORSsettings(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "https://simplegame.now.sh")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, User-Agent, Cache-Control, Accept, X-Requested-With, If-Modified-Since")
		next.ServeHTTP(w, r)
	})
}

func main() {

	http.HandleFunc("/signin", CORSsettings(signinHandler))
	http.HandleFunc("/signup", CORSsettings(signupHandler))
	http.HandleFunc("/profile", CORSsettings(profileHandler))
	http.HandleFunc("/leaders", CORSsettings(leadersHandler))

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
		Online, _ := loggedIn(r)
		if Online {
			http.Redirect(w, r, r.Referer(), http.StatusBadRequest)
			return
		} else {
			http.Redirect(w, r, r.Referer(), http.StatusOK)
			return
		}
	}

	user, err := getFormReq(r)

	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusBadRequest)
	}

	if user.Email != "" && user.Password != "" {
		if !checkExist(*user) {

			user.KeyWord = RandStringBytesRmndr()
			id, _ := addUser(*user)
			session := new(http.Cookie)
			session.Name = "session_id"
			session.Value = uidGen()
			session.Expires = time.Now().Add(time.Minute)
			session.HttpOnly = true
			http.SetCookie(w,session)
			sessions[session.Value] = id

			return
		} else {
			http.Redirect(w, r, r.Referer(), http.StatusAlreadyReported)
			return
		}
	} else {
		http.Redirect(w, r, r.Referer(), http.StatusConflict)
		return
	}
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC(), "Request from", r.URL.String())
	fmt.Println("Method", r.Method)
	Online, _ := loggedIn(r)
	if Online {
		http.Redirect(w, r, r.Referer(), http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodGet {
		http.Redirect(w, r, r.Referer(), http.StatusOK)
		return
	}

	user, err := getFormReq(r)
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusBadRequest)
		return
	}

	if validateUser(*user) {
		session := new(http.Cookie)
		session.Name = "session_id"
		session.Value = uidGen()
		session.Expires = time.Now().Add(time.Hour)
		session.HttpOnly = true
		http.SetCookie(w, session)

		sessions[session.Value] = user.Email

		http.Redirect(w, r, "/profile", http.StatusOK)
	} else {
		http.Redirect(w, r, r.Referer(), http.StatusBadRequest)
		return
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC(), "Request from", r.URL.String())
	fmt.Println("Method", r.Method)

	Online, id := loggedIn(r)

	if r.Method == http.MethodGet {

		if !Online {
			http.Redirect(w, r, "/login", http.StatusBadRequest)
			return
		}
		userJson, err := json.Marshal(users[id])
		if err != nil {
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Status-Code", "200")

		w.Write(userJson)
		return
	}

	if err := uploadFileReq(id, r); err != nil {
		http.Redirect(w, r, "/profile", http.StatusBadRequest)
		return
	}

	user := users[id]
	data, err := getFormReq(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	if data.Nick != "" {
		user.Nick = data.Nick
	}
	if data.Password != "" {
		user.Password = data.Password
	}
	users[id] = user

	http.Redirect(w, r, "/profile", http.StatusOK)
	return

}

func loggedIn(r *http.Request) (bool, string){
	val := r.Cookies()

	for i := 0; i < len(val); i++{
		id, ok := sessions[val[i].Value]
		if !ok {
			continue
		} else {
			return true, users[id].Email
		}

	}

	return false, ""
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

func getFormReq(r *http.Request) (*User, error) {
	user := new(User)
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")
	user.Name = r.FormValue("name")
	user.LastName = r.FormValue("last_name")
	user.Nick = r.FormValue("nick")

	return user, nil

	// body, err := ioutil.ReadAll(r.Body)
	// defer r.Body.Close() // важный пункт!
	// if err != nil {
	// 	return nil, err
	// }
	// user := new(User)
	// err = json.Unmarshal(body, user)
	// return user, err
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
