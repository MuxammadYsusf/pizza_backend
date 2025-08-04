package service

import (
	"context"
	"fmt"
	"github/http/copy/task4/generated/session"
	"strings"
)

func (s *LoginService) Register(ctx context.Context, req *session.RegisterRequest) (*session.RegisterResponse, error) {

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("invalid name or password or phone number")
	}

	if strings.ContainsAny(req.Username, "!@#$%&*()№:;<>?'\";--") {
		return nil, fmt.Errorf("invalid name")
	}

	if strings.ContainsAny(req.Password, "'\";--") {
		return nil, fmt.Errorf("invalid password")
	}

	if len(req.Username) > 15 || len(req.Password) > 50 {
		return nil, fmt.Errorf("invalid name or password")
	}

	if req.Password == "admin" {
		req.Role = "admin"
	} else {
		req.Role = "user"
	}

	resp, err := s.loginPostgres.Login().Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *LoginService) Login(ctx context.Context, req *session.LoginRequest) (*session.LoginResponse, error) {

	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("invalid name or password")
	}

	if strings.ContainsAny(req.Username, "!@#$%&*()№:;<>?") {
		return nil, fmt.Errorf("invalid name or password")
	}

	if strings.ContainsAny(req.Password, "'\";--") {
		return nil, fmt.Errorf("invalid name or password")
	}

	if len(req.Username) > 15 || len(req.Password) > 50 {
		return nil, fmt.Errorf("invalid name or password")
	}

	resp, err := s.loginPostgres.Login().Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
