package profile

import (
	//"SimpleGame/2018_2_Simple_Name/internal/dataParsing"
	//"SimpleGame/2018_2_Simple_Name/internal/db/postgres"
	//"SimpleGame/2018_2_Simple_Name/internal/session"
	"SimpleGame/internal/dataParsing"
	"SimpleGame/internal/db/postgres"
	"SimpleGame/internal/session"

	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var logger, _ = zap.NewProduction()
var sugar = logger.Sugar()

func ProfileHandler(w http.ResponseWriter, r *http.Request) { // Валидировать данные

	sess, err := session.FindSession(r)

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
		fmt.Println(sess.Email)
		user, err := db.GetUser(sess.Email)

		if err != nil { // Полная проверка ошибки?
			//fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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

		existUserData, err := db.GetUser(sess.Email)

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

		newUserData, err := dataParsing.GetJSONReq(r)

		if err != nil {
			sugar.Errorw("Failed get json",
				"error", err,
				"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		validData := govalidator.HasUpperCase(newUserData.Password) && govalidator.HasLowerCase(newUserData.Password)

		if validData {
			user, err := db.UpdateUser(existUserData, newUserData)

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
		if err := UploadFileReq(sess.Email, r); err != nil {
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

func UploadFileReq(fileName string, r *http.Request) error {
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

func GetAvatar(w http.ResponseWriter, r *http.Request) {
	sess, err := session.FindSession(r)

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