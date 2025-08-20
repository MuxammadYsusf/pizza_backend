package handler

import (
	pb "github/http/copy/task4/genproto"
	"github/http/copy/task4/pkg/security"
	"github/http/copy/task4/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(ctx *gin.Context) {
	var reg pb.RegisterRequest

	if err := ctx.ShouldBindJSON(&reg); err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Auth().Register(ctx, &reg)
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)

}

func (h *Handler) Login(ctx *gin.Context) {
	var l pb.LoginRequest

	if err := ctx.ShouldBindJSON(&l); err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Auth().Login(ctx, &l)
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}
	tokenStr, err := security.GenerateJWT(map[string]interface{}{
		"user_id": resp.UserId,
		"role":    resp.Role,
	}, time.Hour*24, "secret")
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, gin.H{
		"token": tokenStr,
	})
}
