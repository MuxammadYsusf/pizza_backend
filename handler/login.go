package handler

import (
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/generated/session"
	"github/http/copy/task4/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(ctx *gin.Context) {
	var reg session.RegisterRequest

	if err := ctx.ShouldBindJSON(&reg); err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err})
		return
	}

	resp, err := h.GRPCClient.Login().Register(ctx, &reg)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err})
		return
	}

	ctx.JSON(c.OK, resp)

}

func (h *Handler) Login(ctx *gin.Context) {
	var l session.LoginRequest

	if err := ctx.ShouldBindJSON(&l); err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err})
		return
	}

	resp, err := h.GRPCClient.Login().Login(ctx, &l)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err})
		return
	}

	tokenStr, err := security.GenerateJWTToken(int(resp.UserId), resp.Role)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, &gin.H{
		"token": tokenStr,
	})
}
