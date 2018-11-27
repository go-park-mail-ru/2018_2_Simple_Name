package dataParsing

import (
	//"SimpleGame/2018_2_Simple_Name/internal/models"
	"SimpleGame/internal/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetJSONReq(r *http.Request) (*models.User, error) {
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

