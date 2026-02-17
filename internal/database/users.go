package database

import "database/sql"

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID int `json:"id"` 
}

