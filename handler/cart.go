package handler

import (
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/generated/pizza"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PutPizzaIntoCart(ctx *gin.Context) {

	var req struct {
		Items  []*pizza.CartItems `json:"items"`
		UserId int                `json:"userId"`
		Id     int                `json:"id"`
	}

	req.UserId = ctx.GetInt("userId")
	if req.UserId == 0 {
		ctx.JSON(c.UnAuth, gin.H{"error": "Не авторизован"})
		return
	}

	IdStr := ctx.Param("id")

	id, err := strconv.Atoi(IdStr)
	if err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}

	req.Id = id

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(c.Empty, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Cart().Cart(ctx, &pizza.CartRequest{
		Items:  req.Items,
		UserId: int32(req.UserId),
		Id:     int32(req.Id),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) DecreasePizzaQuantity(ctx *gin.Context) {
	var req struct {
		PizzaId  int32 `json:"pizzaId"`
		Quantity int32 `json:"quantity"`
		UserId   int   `json:"userId"`
	}

	req.UserId = ctx.GetInt("userId")
	if req.UserId == 0 {
		ctx.JSON(c.UnAuth, gin.H{"error": "Не авторизован"})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(c.Empty, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Cart().DecreasePizzaQuantity(ctx, &pizza.CartItems{
		UserId:   int32(req.UserId),
		PizzaId:  req.PizzaId,
		Quantity: req.Quantity,
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) GetTotalCost(ctx *gin.Context) {

	IdStr := ctx.Param("id")

	id, err := strconv.Atoi(IdStr)
	if err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Cart().GetTotalCost(ctx, &pizza.CartItems{
		CartId: int32(id),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) GetCartHistory(ctx *gin.Context) {

	userId := ctx.GetInt("userId")
	if userId == 0 {
		ctx.JSON(c.UnAuth, gin.H{"error": "Не авторизован"})
		return
	}

	resp, err := h.GRPCClient.Cart().GetCartHistory(ctx, &pizza.GetCartHistoryRequest{
		UserId: int32(userId),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) GetCartItemHistory(ctx *gin.Context) {

	id := ctx.GetInt("id")

	resp, err := h.GRPCClient.Cart().GetCartItemHistory(ctx, &pizza.GetCarItemtHistoryRequest{
		CartHistoryId: int32(id),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)

}
