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
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
var gameService = game.NewGame()

func main() {

	defer logger.Sync()

	err := postgres.InitService()

	if err != nil {
		sugar.Errorw("Failed connect to the database",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		return
	}

	_, err = redis.InitService()

	if err != nil {
		//logging.ErrorLog("Failed open redis", err, sugar)
		sugar.Errorw("Failed open redis",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		////fmt.Println(err.Error())
		return
	}

	//defer obj.Close() // Не будет работать

	go gameService.Run()

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/signup", CORSsettings(signupHandler))
	siteMux.HandleFunc("/signin", CORSsettings(signinHandler))
	siteMux.HandleFunc("/profile", CORSsettings(profileHandler))
	siteMux.HandleFunc("/leaders", CORSsettings(leadersHandler))
	siteMux.HandleFunc("/islogged", CORSsettings(islogged))
	siteMux.HandleFunc("/logout", CORSsettings(logOut))
	siteMux.HandleFunc("/startgame", startGame)
	siteMux.HandleFunc("/leaderscount", CORSsettings(leadersCount))
	siteMux.HandleFunc("/getAvatar", CORSsettings(getAvatar))

	siteHandler := logging.AccessLogMiddleware(siteMux, sugar)

	port := "8080"

	sugar.Infow("starting server at :" + port)

	//fmt.Println("starting server at :8080")
	if err := http.ListenAndServe(":"+port, siteHandler); err != nil {
		log.Fatalf("cannot listen: %s", err)
	}
}

func startGame(w http.ResponseWriter, r *http.Request) {
	sugar.Info("Startgame signal from user")
	sess, err := findSession(r)
	if err != nil {
		sugar.Errorw("Failed get SESSION",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sess == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := postgres.GetUser(sess.Email)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		origin := r.Header["Origin"]
		if len(origin) == 0 {
			return true
		}
		u, err := url.Parse(origin[0])
		if err != nil {
			return false
		}
		originUrl := "127.0.0.1:3000"
		return strings.EqualFold(u.Host, originUrl)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		sugar.Errorw("Cannot upgrade connection", "Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player := game.NewPlayer(user.Nick, conn)

	gameService.Connection <- player
}

func leadersCount(w http.ResponseWriter, r *http.Request) {
	limit := "50" // Общий лимит на показ лидеров
	count, err := postgres.GetLeadersCount(limit)

	if err != nil {
		sugar.Errorw("Failed get count leaders",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	info := models.DBinfo{}

	info.LeadersCount = count

	if err != nil {
		sugar.Errorw("Failed set json",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(info)

	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Status-Code", "200")

	w.Write(resp)

	return
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
	validPassword := govalidator.HasUpperCase(user.Password) && govalidator.HasLowerCase(user.Password) //&& govalidator.IsByteLength(user.Password, 6, 12)
	validNick := !govalidator.HasWhitespace(user.Nick)

	if validEmail && validPassword && validNick {
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
		w.WriteHeader(http.StatusBadRequest)
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

	if existUser == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil { // Проверка, что ошибка НЕ no rows in result set

		sugar.Errorw("Failed get USER",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
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

func getAvatar(w http.ResponseWriter, r *http.Request) {
	sess, err := findSession(r)

	if err != nil {
		sugar.Errorw("Failed get SESSION",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sess == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	file, err := os.Open("./media/" + sess.Email)

	//res, _ := ioutil.ReadAll(file)

	defer file.Close()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	file.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := file.Stat()                         //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+sess.Email)
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	file.Seek(0, 0)
	io.Copy(w, file) //'Copy' the file to the client
	return
}

func profileHandler(w http.ResponseWriter, r *http.Request) { // Валидировать данные

	vars := mux.Vars(r)

	pId := vars["id"]

	fmt.Println(pId)

	sess, err := findSession(r)

	if err != nil {
		sugar.Errorw("Failed get SESSION",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sess == nil {
		w.WriteHeader(http.StatusUnauthorized)
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

		existUserData, err := postgres.GetUser(sess.Email)

		if existUserData == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err != nil {
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

		validData := govalidator.HasUpperCase(newUserData.Password) && govalidator.HasLowerCase(newUserData.Password)

		if validData {
			user, err := postgres.UpdateUser(existUserData, newUserData)

			if err != nil {
				sugar.Errorw("Failed update user",
					"error", err,
					"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			resp, _ := json.Marshal(user)

			w.Write(resp)

			w.WriteHeader(http.StatusOK)
			return

		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if r.Method == http.MethodPost {
		if err := uploadFileReq(sess.Email, r); err != nil {
			sugar.Errorw("Failed put file",
				"error", err,
				"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return

}

func findSession(r *http.Request) (*models.UserSession, error) {
	val := r.Cookies()

	for i := 0; i < len(val); i++ {
		//fmt.Println(val[i].Value)
		if val[i].Name == "session_id" {
			sess, err := session.Get(redis, val[i].Value)

			if err != nil {
				return nil, err
			}
			return sess, nil

		} else {
			continue
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
	}

	w.WriteHeader(http.StatusUnauthorized)
	return

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
}

func logOut(w http.ResponseWriter, r *http.Request) {
	sess, err := findSession(r)

	if err != nil {
		sugar.Errorw("Failed find SESSION",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sess == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = session.Delete(redis, sess.Id, &w)

	if err != nil {
		sugar.Errorw("Failed delete SESSION",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

// func getFormReq(r *http.Request) (*models.User, error) {
// 	user := new(models.User)
// 	user.Email = r.FormValue("email")
// 	user.Password = r.FormValue("password")
// 	user.Nick = r.FormValue("nick")

// 	return user, nil
// }

func uploadFileReq(fileName string, r *http.Request) error {
	if err := r.ParseMultipartForm(32 << 20); nil != err {
		fmt.Println("3")

		return err
	}

	tt, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	fmt.Println()
	fmt.Println(tt)
	fmt.Println()

	file, _, err := r.FormFile("new_avatar")

	if err != nil {
		fmt.Println("1")
		return err
	}
	defer file.Close()

	fmt.Println(fileName)
	fmt.Println(filepath.Join(("/media")))

	dst, err := os.Create(filepath.Join("./media", fileName))

	if err != nil {
		fmt.Println("2")

		return err
	}

	io.Copy(dst, file)
	return nil
}
