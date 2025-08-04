package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github/http/copy/task4/generated/session"
	"github/http/copy/task4/models"

	"golang.org/x/crypto/bcrypt"
)

type auth struct {
	db *sql.DB
}

type AuthRepo interface {
	Register(ctx context.Context, req *session.RegisterRequest) (*session.RegisterResponse, error)
	Login(ctx context.Context, req *session.LoginRequest) (*session.LoginResponse, error)
}

func NewAuth(db *sql.DB) AuthRepo {
	return &auth{
		db: db,
	}
}

func (a *auth) IsNameTaken(name string) (bool, error) {

	var exists bool
	err := a.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (a *auth) IsEmailTaken(name string) (bool, error) {

	var exists bool
	err := a.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (a *auth) Register(ctx context.Context, req *session.RegisterRequest) (*session.RegisterResponse, error) {

	tx, err := a.db.Begin()
	if err != nil {
		return nil, err
	}

	exists, err := a.IsNameTaken(req.Username)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if exists {
		tx.Rollback()
		return nil, fmt.Errorf("user already exists")
	}

	exists, err = a.IsEmailTaken(req.Email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if exists {
		tx.Rollback()
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	req.Password = hashedPassword

	defer tx.Commit()

	query := `INSERT INTO users(name, password, email, role)
	VALUES($1, $2, $3, $4)`

	_, err = tx.Exec(
		query,
		req.Username,
		hashedPassword,
		req.Email,
		req.Role,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &session.RegisterResponse{
		Message: "success",
	}, nil
}

func (a *auth) Login(ctx context.Context, req *session.LoginRequest) (*session.LoginResponse, error) {

	var user models.User

	query := `SELECT id, name, password, role FROM users WHERE name = $1 AND password = $2`

	err := a.db.QueryRow(query, req.Username, req.Password).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
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
