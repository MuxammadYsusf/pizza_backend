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
