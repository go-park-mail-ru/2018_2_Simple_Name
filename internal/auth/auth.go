package auth

import (
	//"SimpleGame/2018_2_Simple_Name/internal/dataParsing"
	//"SimpleGame/2018_2_Simple_Name/internal/db/postgres"
	//"SimpleGame/2018_2_Simple_Name/internal/session"
	//"SimpleGame/2018_2_Simple_Name/internal/validation"
	"SimpleGame/internal/dataParsing"
	"SimpleGame/internal/db/postgres"
	"SimpleGame/internal/session"
	"SimpleGame/internal/validation"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction()
var sugar = logger.Sugar()

//var postgres models.UserService = &db.PostgresUserService{}

func SignupHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//user, err := getFormReq(r)

	user, err := dataParsing.GetJSONReq(r)

	if err != nil {
		sugar.Errorw("Failed get JSON",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	existUser, err := db.GetUser(user.Email)

	if err != nil {
		//fmt.Println("Getuser error: ", err.Error())
	}

	validUser := validation.ValidUser(user)

	if validUser {
		if existUser == nil {

			err := db.CreateUser(user)

			if err != nil {
				sugar.Errorw("Failed create USER",
					"error", err,
					"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			//err = session.Create(redis, user, &w)
			err = session.SessionObj.SetCookie(user, &w)

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

func SigninHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	//user, err := getFormReq(r)

	user, err := dataParsing.GetJSONReq(r)

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

	existUser, err := db.GetUser(user.Email)

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

		//session.Create(redis, user, &w)
		err := session.SessionObj.SetCookie(user, &w)

		if err != nil {
			sugar.Errorw("Failed create SESSION",
				"error", err,
				"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)

		return
	} else {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	//sess, err := findSession(r)
	err := session.SessionObj.RmCookie(r, &w)

	if err != nil {
		sugar.Errorw("Failed find SESSION",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//if sess == nil {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	//_, err = session.Delete(redis, sess.Id, &w)

	//if err != nil {
	//	sugar.Errorw("Failed delete SESSION",
	//		"error", err,
	//		"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	return
}

func Islogged(w http.ResponseWriter, r *http.Request) {

	sess, err := session.SessionObj.FindSession(r)

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
