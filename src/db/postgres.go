package db

import (
	"SimpleGame/models"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresUserService struct {
	db *sql.DB
}

func (s *PostgresUserService) InitService() error {
	var err error
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	s.db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		return err
	}

	if err = s.db.Ping(); err != nil {
		return err
	}

	fmt.Println("You connected to the database.")

	return nil
}


func (s *PostgresUserService) GetUser(email string) (*models.User, error) {
	user := new(models.User)

	query := "SELECT * FROM users WHERE email=$1"
	row := s.db.QueryRow(query, email)

	err := row.Scan(&user.Email, &user.Name, &user.LastName, &user.Age, &user.Nick, &user.Password, &user.Score)

	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s *PostgresUserService) CreateUser(u *models.User) (error) {
	query := "INSERT INTO users(email, name, last_name, age, nick, password) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := s.db.Exec(query,
		u.Email, u.Name, u.LastName, u.Age, u.Nick, u.Password)

	if err != nil {
		return err
	}

	return nil
}
func (s *PostgresUserService) DeleteUser(email string) (error) {
	query := "DELETE FROM users WHERE email=$1"
	_, err := s.db.Exec(query, email)

	if err != nil {
		return err
	}

	return nil
}
func (s *PostgresUserService) UpdateUser(existData *models.User, newData *models.User) (*models.User, error) {
	query := "UPDATE users SET password=$1 WHERE email=$2 RETURNING *"

	row := s.db.QueryRow(query, newData.Password, existData.Email)

	user := new(models.User)

	err := row.Scan(&user.Email, &user.Name, &user.LastName, &user.Age, &user.Nick, &user.Password, &user.Score)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserService) GetUsersByScore(limit string, offset string) ([]*models.User, error) {
	var users = make([]*models.User, 0)

	query := "SELECT nick, score FROM users ORDER BY score DESC LIMIT "+ limit + " OFFSET " + offset + ";"

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next(){
		user := new(models.User)

		err := rows.Scan(&user.Nick, &user.Score)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil

}

func (s *PostgresUserService) GetLeadersCount(limit string) (int, error) {
	query := "SELECT COUNT(*) FROM (SELECT * FROM users ORDER BY score LIMIT "+ limit + ") as foo;"

	row := s.db.QueryRow(query)

	var count int

	err := row.Scan(&count)

	if err != nil {
		return -1, err
	}

	return count, nil

}

