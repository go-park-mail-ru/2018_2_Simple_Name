package leaderboard

import (
	//"SimpleGame/2018_2_Simple_Name/internal/db/postgres"
	//"SimpleGame/2018_2_Simple_Name/internal/models"
	"SimpleGame/internal/models"
	"SimpleGame/internal/db/postgres"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
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

	resp, _ := json.Marshal(info)

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

	top, err := db.GetUsersByScore(limit, offset)

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
