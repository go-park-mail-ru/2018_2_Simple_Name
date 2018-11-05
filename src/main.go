package main

import (
	"./db"
	"./models"
	"./session"
	"SimpleGame/2018_2_Simple_Name/src/logging"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

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

func main() {

	defer logger.Sync()


	postgres.InitService()
	obj, err := redis.InitService()

	if err != nil {
		//logging.ErrorLog("Failed open redis", err, sugar)
		sugar.Errorw("Failed open redis",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
		////fmt.Println(err.Error())
		return
	}

	defer obj.Close() // Не будет работать

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/sugnup", CORSsettings(signupHandler))
	siteMux.HandleFunc("/signup", CORSsettings(signupHandler))
	siteMux.HandleFunc("/profile", CORSsettings(profileHandler))
	siteMux.HandleFunc("/leaders", CORSsettings(leadersHandler))
	siteMux.HandleFunc("/islogged", CORSsettings(islogged))

	siteHandler := logging.AccessLogMiddleware(siteMux, sugar)

	sugar.Infow("starting server at :8080")

	//fmt.Println("starting server at :8080")
	if err := http.ListenAndServe(":8080", siteHandler); err != nil {
		log.Fatalf("cannot listen: %s", err)
	}
}

func leadersHandler(w http.ResponseWriter, r *http.Request) {

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

func ValidUser(user *models.User) (bool, error) {
	validEmail := govalidator.IsEmail(user.Email)
	validPassword := govalidator.HasUpperCase(user.Password) && govalidator.HasLowerCase(user.Password) && govalidator.IsByteLength(user.Password, 6,12)
	validNick := !govalidator.HasWhitespace(user.Nick)
	validName := govalidator.IsAlpha(user.Name)
	validLastName := govalidator.IsAlpha(user.LastName)

	if validEmail && validPassword && validNick && validName && validLastName{
		return true, nil
	} else {
		return false, nil
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
			"time", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	existUser, err := postgres.GetUser(user.Email)

	if err != nil {
		//fmt.Println("Getuser error: ", err.Error())
	}

	validUser, _ := ValidUser(user)

	if validUser {
		fmt.Println("ВАЛИД")
	} else {
		fmt.Println("НЕ ВАЛИД")
	}

	if validUser {
		if existUser == nil {

			err := postgres.CreateUser(user)

			if err != nil {
				sugar.Errorw("Failed create USER",
					"error", err,
					"time", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = session.Create(redis, user, &w)

			if err != nil {
				sugar.Errorw("Failed create SESSION",
					"error", err,
					"time", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
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
			"time", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
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

func profileHandler(w http.ResponseWriter, r *http.Request) {// Валидировать данные

	sess, err := findSession(r)

	if err != nil || sess == nil {
		sugar.Errorw("Failed get SESSION",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
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
				"time", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
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
				"time", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
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
		//fmt.Println(err.Error())
		return err
	}

	file, _, err := r.FormFile("my_file")
	if err != nil {
		//fmt.Println(err.Error())
		return nil
	}
	defer file.Close()

	dst, err1 := os.Create(filepath.Join("../src/static/media", fileName))

	if err1 != nil {
		//fmt.Println(err1.Error())
		return err
	}

	io.Copy(dst, file)
	return nil
}
