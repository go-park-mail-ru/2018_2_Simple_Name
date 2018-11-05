package main

import (
	"./models"
	"./db"
	"./session"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)


//type Sessions map[string]string

//var sessions = make(Sessions)

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

var postgres models.UserService = &db.PostgresUserService{}
var redis models.UserSessionService = &session.RedisSessionService{}

func main() {
	postgres.InitService()
	obj, err := redis.InitService()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer obj.Close() // Не будет работать

	//var session = models.UserSession{}

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

	Leaders := map[int]models.User{
		0: models.User{
			Nick:  "GRe12",
			Score: 4321,
			Age:   12,
		},
		1: models.User{
			Nick:  "wasaW2",
			Score: 43121,
			Age:   13,
		},
		2: models.User{
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

	existUser, err := postgres.GetUser(user.Email)

	if err != nil {
		fmt.Println("Getuser error: ", err.Error())
	}

	if user.Email != "" && user.Password != "" {
		if existUser == nil {

			err := postgres.CreateUser(user)

			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = session.Create(redis, user, &w)

			if err != nil {
				fmt.Println("Ошибка ses.create: ", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

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

	existUser, err := postgres.GetUser(user.Email)

	if err != nil || existUser == nil {
		fmt.Println("Signin fail: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if existUser.Password == user.Password {

		session.Create(redis, user, &w)

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

	sess, err := findSession(r)

	if err != nil || sess == nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		user, err := postgres.GetUser(sess.Email)

		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userInfo, err := json.Marshal(user)

		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Status-Code", "200")

		w.Write(userInfo)
		return
	}

	//if err := uploadFileReq(id, r); err != nil {
	//
	//	return
	//}

	//user := users[id]
	//data, err := getFormReq(r)
	//if err != nil {
	//	return
	//}
	//if data.Nick != "" {
	//	user.Nick = data.Nick
	//}
	//if data.Password != "" {
	//	user.Password = data.Password
	//}
	//users[id] = user

	return

}

func findSession(r *http.Request) (*models.UserSession, error) {
	val := r.Cookies()

	for i := 0; i < len(val); i++{
		fmt.Println(val[i].Value)
		sess, err := session.Get(redis, val[i].Value)
		if sess == nil {
			continue
		} else {
			return sess, nil
		}
		if err != nil {
			fmt.Println("islogged error: ", err.Error())
			return nil, err
		}

	}
	return nil, nil
}

func islogged(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC(), "Request from", r.URL.String())
	fmt.Println("Method", r.Method)

	sess, err := findSession(r)

	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sess != nil {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//val := r.Cookies()
	//
	//for i := 0; i < len(val); i++{
	//	fmt.Println(val[i].Value)
	//	sess, err := session.Get(redis, val[i].Value)
	//	if sess == nil {
	//		continue
	//	}
	//	if err != nil {
	//		fmt.Println("islogged error: ", err.Error())
	//		return
	//	}
	//
	//	if sess.Email == "" { // != nil? можно поставить так.
	//		continue
	//	} else {
	//		w.WriteHeader(http.StatusOK)
	//		return
	//	}
	//
	//}
	w.WriteHeader(http.StatusBadRequest)
	return
}

func getJSONReq(r *http.Request) (*models.User, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		fmt.Println("Ошибка чтения 1: ", err.Error())
		return nil, err
	}

	user := new(models.User)

	err = json.Unmarshal(body, user)

	if err != nil {
		fmt.Println("Ошибка чтения 2: ", err.Error())
		return nil, err
	}

	return user, nil
}

func getFormReq(r *http.Request) (*models.User, error) {
	user := new(models.User)
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
