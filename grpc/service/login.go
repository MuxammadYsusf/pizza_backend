package service

import (
	"context"
	"fmt"
	"github/http/copy/task4/generated/session"
	"github/http/copy/task4/pkg/helper"
	"github/http/copy/task4/security"
	"strings"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AuthService) Register(ctx context.Context, req *session.RegisterRequest) (*session.RegisterResponse, error) {

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("invalid name or password or email address")
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

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashedPassword

	resp, err := s.authPostgres.Auth().Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

var counter int64

func tokenId() int {

	id := atomic.AddInt64(&counter, 1)

	return int(id)
}

func (s *AuthService) Login(ctx context.Context, req *session.LoginRequest) (*session.LoginResponse, error) {

	id := tokenId()

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

	resp, err := s.authPostgres.Auth().GetUserData(ctx, req)
	if err != nil {
		return nil, err
	}

	data, err := s.authPostgres.Auth().GetSessionData(ctx, id)
	if err != nil {
		return nil, err
	}

	si, err := s.authPostgres.Auth().GetSessionId(ctx)
	if err != nil {
		return nil, err
	}

	if si.ID > 0 && time.Now().Before(data.ExpiredAt) {
		id = si.ID
	}

	id = si.ID + 1

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if !helper.CheckPasswordHash(req.Password, hashedPassword) {
		return nil, nil
	}

	tokenStr, err := security.GenerateJWTToken(int(resp.Id), id, resp.Role)
	if err != nil {
		return nil, err
	}

	createdAt := time.Now()
	expiredAt := createdAt.Add(time.Hour * 24)

	_, err = s.authPostgres.Auth().InsertSession(ctx, &session.InsertSessionRequest{
		Id:        int32(id),
		UserId:    int32(resp.Id),
		CreatedAt: timestamppb.New(createdAt),
		ExpiredAt: timestamppb.New(expiredAt),
	})
	if err != nil {
		return nil, err
	}

	resp = &session.LoginResponse{
		Token: tokenStr,
	}

	return resp, nil
}

func (s *AuthService) Logout(ctx context.Context, req *session.LogoutRequest) (*session.LogoutResponse, error) {

	si, err := s.authPostgres.Auth().GetSessionId(ctx)
	if err != nil {
		return nil, err
	}

	req.Id = int32(si.ID)

	resp, err := s.authPostgres.Auth().Logout(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
