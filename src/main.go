package main

import (
	"SimpleGame/db"
	"SimpleGame/game"
	"SimpleGame/logging"
	"SimpleGame/models"
	"SimpleGame/session"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func CORSsettings(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, User-Agent, Cache-Control, Accept, X-Requested-With, If-Modified-Since")
		if r.Method == http.MethodOptions {
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func IsLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := findSession(r)

		if err != nil {
			//fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if sess != nil {
			//w.WriteHeader(http.StatusOK)
			//return

			next.ServeHTTP(w, r)
		} else {
			//w.WriteHeader(http.StatusBadRequest)
			next.ServeHTTP(w, r)
			return
		}
	})

}

var postgres models.UserService = &db.PostgresUserService{}
var redis models.UserSessionService = &session.RedisSessionService{}
var logger, _ = zap.NewProduction()
var sugar = logger.Sugar()
var g *game.Game = game.NewGame()

func main() {

	defer logger.Sync()

	postgres.InitService()
	obj, err := redis.InitService()

	if err != nil {
		//logging.ErrorLog("Failed open redis", err, sugar)
		sugar.Errorw("Failed open redis",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		////fmt.Println(err.Error())
		return
	}

	defer obj.Close() // Не будет работать

	go g.Run()

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/signup", CORSsettings(signupHandler))
	siteMux.HandleFunc("/signin", CORSsettings(signinHandler))
	siteMux.HandleFunc("/profile", CORSsettings(profileHandler))
	siteMux.HandleFunc("/leaders", CORSsettings(leadersHandler))
	siteMux.HandleFunc("/islogged", CORSsettings(islogged))
	siteMux.HandleFunc("/startgame", CORSsettings(startGame))

	siteHandler := logging.AccessLogMiddleware(siteMux, sugar)

	sugar.Infow("starting server at :8080")

	//fmt.Println("starting server at :8080")
	if err := http.ListenAndServe(":8080", siteHandler); err != nil {
		log.Fatalf("cannot listen: %s", err)
	}
}

func startGame(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		sugar.Errorw("Cannot upgrade connection", "Error:",err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nickname:="nick"///////////////////fix it

	g.Connection <- &game.Player{Conn:conn,Nickname:nickname}
}

func leadersHandler(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	_, err := strconv.Atoi(limit)

	if err != nil {
		limit = "ALL"
	}

	_, err = strconv.Atoi(offset)

	if err != nil {
		offset = "0"
	}

	top, err := postgres.GetUsersByScore(limit, offset)

	if err != nil {
		sugar.Errorw("Failed get users",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(top)

	if err != nil {
		sugar.Errorw("Failed set JSON",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Status-Code", "200")

	w.Write(resp)

	return
}
func ValidUser(user *models.User) bool {
	validEmail := govalidator.IsEmail(user.Email)
	validPassword := govalidator.HasUpperCase(user.Password) && govalidator.HasLowerCase(user.Password) && govalidator.IsByteLength(user.Password, 6, 12)
	validNick := !govalidator.HasWhitespace(user.Nick)
	validName := govalidator.IsAlpha(user.Name) && !govalidator.HasWhitespace(user.Nick)
	validLastName := govalidator.IsAlpha(user.LastName) && !govalidator.HasWhitespace(user.Nick)

	if validEmail && validPassword && validNick && validName && validLastName {
		return true
	} else {
		return false
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//user, err := getFormReq(r)

	user, err := getJSONReq(r)

	if err != nil {
		sugar.Errorw("Failed get JSON",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	existUser, err := postgres.GetUser(user.Email)

	if err != nil {
		//fmt.Println("Getuser error: ", err.Error())
	}

	validUser := ValidUser(user)

	if validUser {
		if existUser == nil {

			err := postgres.CreateUser(user)

			if err != nil {
				sugar.Errorw("Failed create USER",
					"error", err,
					"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = session.Create(redis, user, &w)

			if err != nil {
				sugar.Errorw("Failed create SESSION",
					"error", err,
					"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
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

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//user, err := getFormReq(r)

	user, err := getJSONReq(r)

	if err != nil {
		sugar.Errorw("Failed get JSON",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validEmail := govalidator.IsEmail(user.Email)

	if !validEmail {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existUser, err := postgres.GetUser(user.Email)

	//if err != nil || existUser == nil {
	//	sugar.Errorw("Failed get USER",
	//		"error", err,
	//		"time", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	if existUser.Password == user.Password {

		session.Create(redis, user, &w)

		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) { // Валидировать данные

	vars := mux.Vars(r)

	pId := vars["id"]

	fmt.Println(pId)

	sess, err := findSession(r)

	if err != nil || sess == nil {
		sugar.Errorw("Failed get SESSION",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		user, err := postgres.GetUser(sess.Email)

		//if err != nil { // Полная проверка ошибки?
		//	//fmt.Println(err.Error())
		//	w.WriteHeader(http.StatusInternalServerError)
		//	return
		//}

		userInfo, err := json.Marshal(user)

		if err != nil {
			sugar.Errorw("Failed marshal json",
				"error", err,
				"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Status-Code", "200")

		w.Write(userInfo)
		return
	} else if r.Method == http.MethodPut {

		if err := uploadFileReq(pId, r); err != nil {
			sugar.Errorw("Failed put file",
				"error", err,
				"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		existUserData, err := postgres.GetUser(sess.Email)

		if err != nil || existUserData == nil {
			sugar.Errorw("Failed get user",
				"error", err,
				"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newUserData, err := getJSONReq(r)

		if err != nil {
			sugar.Errorw("Failed get json",
				"error", err,
				"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		validData := ValidUser(newUserData)

		if validData {
			err := postgres.UpdateUser(existUserData, newUserData)

			if err != nil {
				sugar.Errorw("Failed update user",
					"error", err,
					"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			return

		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return

}

func findSession(r *http.Request) (*models.UserSession, error) {
	val := r.Cookies()

	for i := 0; i < len(val); i++ {
		//fmt.Println(val[i].Value)
		sess, err := session.Get(redis, val[i].Value)

		if err != nil {
			return nil, err
		}

		if sess == nil {
			continue
		} else {
			return sess, nil
		}

	}
	return nil, nil
}

func islogged(w http.ResponseWriter, r *http.Request) {

	sess, err := findSession(r)

	if err != nil {
		sugar.Errorw("Failed find SESSION",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
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
	////	fmt.Println(val[i].Value)
	//	sess, err := session.Get(redis, val[i].Value)
	//	if sess == nil {
	//		continue
	//	}
	//	if err != nil {
	////		fmt.Println("islogged error: ", err.Error())
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
		//fmt.Println("Ошибка чтения 1: ", err.Error())
		return nil, err
	}

	user := new(models.User)

	err = json.Unmarshal(body, user)

	if err != nil {
		//fmt.Println("Ошибка чтения 2: ", err.Error())
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
		return err
	}

	file, _, err := r.FormFile("my_file")
	if err != nil {
		return nil
	}
	defer file.Close()

	dst, err1 := os.Create(filepath.Join("../src/static/media", fileName))

	if err1 != nil {
		return err
	}

	io.Copy(dst, file)
	return nil
}
