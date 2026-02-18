package database

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// SQLite doesn't support RETURNING, use LastInsertId
	query := "INSERT INTO users (email, password, name) VALUES (?, ?, ?)"
	res, err := m.DB.ExecContext(ctx, query, user.Email, user.Password, user.Name)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.Id = int(id)
	return nil
}

func (m *UserModel) Get(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, email, name, password FROM users WHERE id = ?"

	var user User
	err := m.DB.QueryRowContext(ctx, query, id).
		Scan(&user.Id, &user.Email, &user.Name, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail finds a user by email (useful for login)
func (m *UserModel) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, email, name, password FROM users WHERE email = ?"

	var user User
	err := m.DB.QueryRowContext(ctx, query, email).
		Scan(&user.Id, &user.Email, &user.Name, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

