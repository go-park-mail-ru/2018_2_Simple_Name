package leaderboard

import (
	//"SimpleGame/2018_2_Simple_Name/internal/db/postgres"
	//"SimpleGame/2018_2_Simple_Name/internal/models"
	"SimpleGame/internal/db/postgres"
	"SimpleGame/internal/models"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction()
var sugar = logger.Sugar()

func LeadersCount(w http.ResponseWriter, r *http.Request) {
	limit := "50" // Общий лимит на показ лидеров
	count, err := db.GetLeadersCount(limit)

	if err != nil {
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

	resp, _ := info.MarshalJSON()

	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Status-Code", "200")

	w.Write(resp)

	return
}

func LeadersHandler(w http.ResponseWriter, r *http.Request) {
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

	usrList := models.UserList{}
	usrList, err = db.GetUsersByScore(limit, offset)

	if err != nil {
		sugar.Errorw("Failed get users",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := usrList.MarshalJSON() //easyjson

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
