package models

// easyjson:json
type User struct {
	Nick     string `json:"nick"`
	Email    string `json:"email"`
	Password string `json:"password"`
	KeyWord  string `json:"-"`
	Score    int    `json:"score"`
}

// easyjson:json
type UserList []*User

//type UserSession struct {
//	Id    string
//	Email string
//}

// easyjson:json
type DBinfo struct {
	LeadersCount int `json:"leaderscount"`
}

type UserService interface {
	InitService() error
	GetUser(email string) (*User, error)
	CreateUser(u *User) error
	UpdateUser(existData *User, newData *User) (*User, error)
	DeleteUser(email string) error
	GetUsersByScore(limit string, offset string) ([]*User, error)
	GetLeadersCount(limit string) (int, error)
}

//
//type UserSessionService interface {
//	InitService() (redis.Conn, error)
//	Create(key string, value string) error
//	Get(key string) (string, error)
//	Delete(key string) error
//}

//type UserSessionService interface {
//
//	Create(user *User, w http.ResponseWriter) (error)
//	Delete(user *User) (error)
//	Get(user *User) (*UserSession, error)
//}

//var redis UserSessionService = new(session.RedisSessionService)
