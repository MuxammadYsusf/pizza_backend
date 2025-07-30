package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github/http/copy/task4/generated/session"
	"github/http/copy/task4/models"

	"golang.org/x/crypto/bcrypt"
)

type login struct {
	db *sql.DB
}

type LoginRepo interface {
	Reg(ctx context.Context, req *session.RegisterRequest) (*session.RegisterResponse, error)
	Login(ctx context.Context, req *session.LoginRequest) (*session.LoginResponse, error)
}

func NewLogin(db *sql.DB) LoginRepo {
	return &login{
		db: db,
	}
}

func (l *login) IsNameTaken(name string) (bool, error) {

	var exists bool
	err := l.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (l *login) IsEmailTaken(name string) (bool, error) {

	var exists bool
	err := l.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (l *login) Reg(ctx context.Context, req *session.RegisterRequest) (*session.RegisterResponse, error) {

	tx, err := l.db.Begin()
	if err != nil {
		return nil, err
	}

	exists, err := l.IsNameTaken(req.Username)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if exists {
		tx.Rollback()
		return nil, fmt.Errorf("user already exists")
	}

	exists, err = l.IsNameTaken(req.Email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if exists {
		tx.Rollback()
		return nil, fmt.Errorf("user already exists")
	}

	defer tx.Commit()

	query := `INSERT INTO users VALUES(name, password, email)`

	_, err = tx.Exec(
		query,
		req.Username,
		req.Password,
		req.Email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &session.RegisterResponse{
		Message: "success",
	}, nil
}

func (l *login) Login(ctx context.Context, req *session.LoginRequest) (*session.LoginResponse, error) {

	var user models.User

	query := `SELECT id, name, password, role FROM users WHERE name = $1 AND password = $2`

	err := l.db.QueryRow(query, req.Username, req.Password).Scan(user.ID, user.Username, user.Password, user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid name or password")
	}

	return &session.LoginResponse{
		Message: "success",
		UserId:  int32(user.ID),
	}, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err == nil
	}

	return err == nil
}
