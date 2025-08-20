package handler

import (
	"fmt"
	pb "github/http/copy/task4/genproto"
	"github/http/copy/task4/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PutPizzaIntoCart(ctx *gin.Context) {

	var req struct {
		Items  []*pb.CartItems `json:"items"`
		UserId int             `json:"userId"`
	}

	req.UserId = ctx.GetInt("userId")
	if req.UserId == 0 {
		ctx.JSON(util.HTTPForbidden, gin.H{"error": "Не авторизован"})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Cart().Cart(ctx, &pb.CartRequest{
		Items:  req.Items,
		UserId: int32(req.UserId),
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) DecreasePizzaQuantity(ctx *gin.Context) {
	var req struct {
		PizzaId  int32 `json:"pizzaId"`
		Quantity int32 `json:"quantity"`
		UserId   int   `json:"userId"`
	}

	req.UserId = ctx.GetInt("userId")
	if req.UserId == 0 {
		ctx.JSON(util.HTTPForbidden, gin.H{"error": "Не авторизован"})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Cart().DecreasePizzaQuantity(ctx, &pb.CartItems{
		UserId:   int32(req.UserId),
		PizzaId:  req.PizzaId,
		Quantity: req.Quantity,
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) GetFromCart(ctx *gin.Context) {
	userId := ctx.GetInt("userId")

	resp, err := h.GRPCClient.Cart().GetFromCart(ctx, &pb.CartItems{
		UserId: int32(userId),
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) ClearTheCart(ctx *gin.Context) {
	userId := ctx.GetInt("userId")

	resp, err := h.GRPCClient.Cart().ClearTheCart(ctx, &pb.CartItems{
		UserId: int32(userId),
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) ClearTheCartById(ctx *gin.Context) {
	userId := ctx.GetInt("userId")

	IdStr := ctx.Param("pizzaId")

	id, err := strconv.Atoi(IdStr)
	if err != nil {
		fmt.Println("HIERE")
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.GRPCClient.Cart().ClearTheCartById(ctx, &pb.CartItems{
		UserId:  int32(userId),
		PizzaId: int32(id),
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) GetTotalCost(ctx *gin.Context) {

	IdStr := ctx.Param("id")

	id, err := strconv.Atoi(IdStr)
	if err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Cart().GetTotalCost(ctx, &pb.CartItems{
		CartId: int32(id),
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) GetCartHistory(ctx *gin.Context) {

	userId := ctx.GetInt("userId")
	if userId == 0 {
		ctx.JSON(util.HTTPForbidden, gin.H{"error": "Не авторизован"})
		return
	}

	resp, err := h.GRPCClient.Cart().GetCartHistory(ctx, &pb.GetCartHistoryRequest{
		UserId: int32(userId),
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) GetCartItemHistory(ctx *gin.Context) {

	IdStr := ctx.Param("id")

	id, err := strconv.Atoi(IdStr)
	if err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Cart().GetCartItemHistory(ctx, &pb.GetCartItemHistoryRequest{
		OrderId: int32(id),
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)

}
