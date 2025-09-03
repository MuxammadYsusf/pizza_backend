package handler

import (
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/generated/session"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(ctx *gin.Context) {
	var reg session.RegisterRequest

	if err := ctx.ShouldBindJSON(&reg); err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Auth().Register(ctx, &reg)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)

}

func (h *Handler) Login(ctx *gin.Context) {
	var l session.LoginRequest

	if err := ctx.ShouldBindJSON(&l); err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Auth().Login(ctx, &l)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) Logout(ctx *gin.Context) {

	auth := ctx.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
	if token == "" {
		ctx.JSON(c.UnAuth, gin.H{"error": "unauthorized"})
		return
	}

	req := session.LogoutRequest{
		Token: token,
	}

	resp, err := h.GRPCClient.Auth().Logout(ctx, &req)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)

}

func (h *Handler) GetUserData(ctx *gin.Context) {

	userId := ctx.GetInt("userId")

	resp, err := h.GRPCClient.Auth().GetUserData(ctx, &session.LoginRequest{
		Id: int32(userId),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateUserPassword(ctx *gin.Context) {
	var req struct {
		l      session.UpdatePasswordRequest
		UserId int `json:"userId"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Auth().UpdateUserPassword(ctx, &session.UpdatePasswordRequest{
		UserId:      int32(req.UserId),
		OldPassword: req.l.OldPassword,
		NewPassword: req.l.NewPassword,
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
