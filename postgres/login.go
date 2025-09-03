package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github/http/copy/task4/generated/session"
	"github/http/copy/task4/models"
	"time"
)

type auth struct {
	db *sql.DB
}

type AuthRepo interface {
	Register(ctx context.Context, req *session.RegisterRequest) (*session.RegisterResponse, error)
	GetUserData(ctx context.Context, req *session.LoginRequest) (*session.LoginResponse, error)
	InsertSession(ctx context.Context, req *session.InsertSessionRequest) (*session.InsertSessionResponse, error)
	Logout(ctx context.Context, req *session.LogoutRequest) (*session.LogoutResponse, error)
	IsLogined(ctx context.Context, req *session.LogoutRequest) (*session.LogoutResponse, error)
	GetSessionByID(ctx context.Context, id int) (*models.Session, error)
	GetSessionData(ctx context.Context, id int) (*models.Session, error)
	GetUserName(ctx context.Context, req *session.LoginRequest) (*session.LoginResponse, error)
	UpdatePassword(ctx context.Context, req *session.UpdatePasswordRequest) (*session.UpdatePasswordResponse, error)
	GetSessionId(ctx context.Context) (*models.Session, error)
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

func (a *auth) IsEmailTaken(email string) (bool, error) {

	var exists bool
	err := a.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
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

	query := `INSERT INTO users(name, password, email, role)
	VALUES($1, $2, $3, $4)`

	_, err = tx.Exec(
		query,
		req.Username,
		req.Password,
		req.Email,
		req.Role,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &session.RegisterResponse{
		Message: "success",
	}, nil
}

func (a *auth) GetUserName(ctx context.Context, req *session.LoginRequest) (*session.LoginResponse, error) {

	var (
		user models.User
	)

	query := `SELECT name FROM users WHERE id = $1`

	err := a.db.QueryRow(query, req.Id).Scan(&user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &session.LoginResponse{
		Username: user.Username,
	}, nil
}

func (a *auth) GetUserData(ctx context.Context, req *session.LoginRequest) (*session.LoginResponse, error) {

	var (
		user models.User
	)

	query := `SELECT id, name, password, role FROM users WHERE name = $1`

	err := a.db.QueryRow(query, req.Username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &session.LoginResponse{
		Id:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}, nil
}

func (a *auth) InsertSession(ctx context.Context, req *session.InsertSessionRequest) (*session.InsertSessionResponse, error) {

	query := `INSERT INTO sessions(id, user_id, created_at, expired_at)
	VALUES($1, $2, $3, $4)`

	_, err := a.db.Exec(
		query,
		req.Id,
		req.UserId,
		req.CreatedAt.AsTime(),
		req.ExpiredAt.AsTime(),
	)
	if err != nil {

		return nil, err
	}

	return &session.InsertSessionResponse{
		Message: "success",
	}, nil
}

func (a *auth) Logout(ctx context.Context, req *session.LogoutRequest) (*session.LogoutResponse, error) {

	query := `UPDATE sessions SET expired_at = NOW() WHERE id = $1`

	_, err := a.db.Exec(
		query,
		req.Id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &session.LogoutResponse{
		Message: "success",
	}, nil
}

func (a *auth) IsLogined(ctx context.Context, req *session.LogoutRequest) (*session.LogoutResponse, error) {

	var expiredAt time.Time

	query := `SELECT expired_at FROM sessions WHERE token = $1`

	err := a.db.QueryRow(query, req.Token).Scan(&expiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &session.LogoutResponse{
		Message: "success",
	}, nil
}

func (a *auth) GetSessionByID(ctx context.Context, id int) (*models.Session, error) {

	var session models.Session

	query := `SELECT id, user_id, created_at, expired_at FROM sessions WHERE id = $1`

	err := a.db.QueryRow(query, id).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &session, nil
}

func (a *auth) GetSessionData(ctx context.Context, id int) (*models.Session, error) {

	var session models.Session

	query := `SELECT * FROM sessions WHERE id = $1`

	err := a.db.QueryRow(query, id).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiredAt)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &session, nil
}

func (a *auth) GetSessionId(ctx context.Context) (*models.Session, error) {

	var session models.Session

	query := `SELECT id FROM sessions ORDER BY id DESC LIMIT 1`

	err := a.db.QueryRow(query).Scan(&session.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &models.Session{
		ID: session.ID,
	}, nil
}

func (a *auth) UpdatePassword(ctx context.Context, req *session.UpdatePasswordRequest) (*session.UpdatePasswordResponse, error) {

	query := `UPDATE users SET password = $2 WHERE name = $1`

	_, err := a.db.Exec(query, req.Username, req.NewPassword)
	if err != nil {
		return nil, err
	}

	return &session.UpdatePasswordResponse{
		Message: "success",
	}, nil
}
