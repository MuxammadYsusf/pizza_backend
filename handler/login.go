package handler

import (
	"github/http/copy/task4/generated/session"
	"github/http/copy/task4/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(ctx *gin.Context) {
	var reg session.RegisterRequest

	if err := ctx.ShouldBindJSON(&reg); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	resp, err := h.GRPCClient.Login().Register(ctx, &reg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)

}

func (h *Handler) Login(ctx *gin.Context) {
	var l session.LoginRequest

	if err := ctx.ShouldBindJSON(&l); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	resp, err := h.GRPCClient.Login().Login(ctx, &l)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokenStr, err := security.GenerateJWTToken(int(resp.UserId), resp.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &gin.H{
		"token": tokenStr,
	})
}
