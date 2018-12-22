package db

import (
	//"SimpleGame/2018_2_Simple_Name/internal/models"
	"SimpleGame/internal/models"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var postgres models.UserService = &PostgresUserService{}

var logger, _ = zap.NewProduction()
var sugar = logger.Sugar()

func OpenConn() error {
	err := postgres.InitService()

	if err != nil {
		sugar.Errorw("Failed connect to the database",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		return err
	}

	return nil
}

func GetUser(email string) (*models.User, error) {
	user, err := postgres.GetUser(email)

	if err != nil {
		sugar.Errorw("Failed getUser from DB",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		return nil, err
	}

	return user, nil
}

func CreateUser(u *models.User) error {
	err := postgres.CreateUser(u)

	if err != nil {
		sugar.Errorw("Failed createUser DB",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		return err
	}

	return nil
}

func DeleteUser(email string) error {
	err := postgres.DeleteUser(email)

	if err != nil {
		sugar.Errorw("Failed deleteUser from DB",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		return err
	}

	return nil
}
func UpdateUser(existData *models.User, newData *models.User) (*models.User, error) {
	user, err := postgres.UpdateUser(existData, newData)

	if err != nil {
		sugar.Errorw("Failed updateUser from DB",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		return nil, err
	}

	return user, nil
}
func UpdateScore(nick string, score int) error {
	err := postgres.UpdateScore(nick, score)
	return err
}

func GetUsersByScore(limit string, offset string) ([]*models.User, error) {
	var users = make([]*models.User, 0)
	// usr:= new(models.UserList)
	// usr = postgres.GetUsersByScore(limit, offset)
	users, err := postgres.GetUsersByScore(limit, offset)

	if err != nil {
		sugar.Errorw("Failed GetUsersByScore from DB",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		return nil, err
	}

	return users, nil

}

func GetLeadersCount(limit string) (int, error) {
	count, err := postgres.GetLeadersCount(limit)

	if err != nil {
		sugar.Errorw("Failed GetUsersByScore from DB",
			"error", err,
			"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		return 0, err
	}

	return count, nil

}
