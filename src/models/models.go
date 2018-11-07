package models

import "github.com/gomodule/redigo/redis"

type User struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Nick     string `json:"nick"`
	Email    string `json:"email"`
	Password string `json:"password"`
	KeyWord  string `json:"-"`
	Score    int    `json:"score"`
	Age      int    `json:"age"`
}

type UserSession struct {
	Id string
	Email string
}

type UserService interface {
	InitService() (error)
	GetUser(email string) (*User, error)
	CreateUser(u *User) error
	UpdateUser(existData *User, newData *User) error
	DeleteUser(email string) error
	GetUsersByScore(limit string, offset string) ([]*User, error)
}


type UserSessionService interface {
	InitService() (redis.Conn, error)
	Create(key string, value string) (error)
	Get(key string) (string, error)
//	Delete(user *User) (error)
}

//type UserSessionService interface {
//
//	Create(user *User, w http.ResponseWriter) (error)
//	Delete(user *User) (error)
//	Get(user *User) (*UserSession, error)
//}



//var redis UserSessionService = new(session.RedisSessionService)

