package handler

import (
	"github/http/copy/task4/generated/pizza"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreatePizzaType(ctx *gin.Context) {
	var req pizza.CreatePizzaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	resp, err := h.GRPCClient.Pizza().CreatePizzaType(ctx, &req)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, resp)
}

func (h *Handler) CreatePizza(ctx *gin.Context) {
	var req pizza.CreatePizzaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	resp, err := h.GRPCClient.Pizza().CreatePizza(ctx, &req)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, resp)
}

func (h *Handler) GetPizzas(ctx *gin.Context) {

	resp, err := h.GRPCClient.Pizza().GetPizzas(ctx, &pizza.GetPizzasRequest{})
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, resp)
}

func (h *Handler) GetPizzaById(ctx *gin.Context) {

	id := ctx.GetInt("id")
	typeId := ctx.GetInt("typeId")

	resp, err := h.GRPCClient.Pizza().GetPizzaById(ctx, &pizza.GetPizzaByIdRequest{
		Id:     int32(id),
		TypeId: int32(typeId),
	})
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, resp)
}

func (h *Handler) UpdatePizza(ctx *gin.Context) {

	var req pizza.UpdatePizzaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	resp, err := h.GRPCClient.Pizza().UpdatePizza(ctx, &req)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, resp)
}

func (h *Handler) DeletePizza(ctx *gin.Context) {

	id := ctx.GetInt("id")
	typeId := ctx.GetInt("typeId")

	resp, err := h.GRPCClient.Pizza().DeletePizza(ctx, &pizza.DeletePizzaRequest{
		Id:     int32(id),
		TypeId: int32(typeId),
	})
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, resp)
}

func (h *Handler) PutPizzaIntoCart(ctx *gin.Context) {

	resp, err := h.GRPCClient.Pizza().Cart(ctx, &pizza.CartRequest{})
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, resp)
}

func (h *Handler) UpdatePizzaInCart(ctx *gin.Context) {

	var req pizza.CartItems

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	resp, err := h.GRPCClient.Pizza().UpdatePizzaInCart(ctx, &req)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, resp)
}

func (h *Handler) OrderPizza(ctx *gin.Context) {

	var req pizza.OrderPizzaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	resp, err := h.GRPCClient.Pizza().OrderPizza(ctx, &req)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, resp)
}

func (h *Handler) GetCartHistory(ctx *gin.Context) {

	resp, err := h.GRPCClient.Pizza().GetUserHistory(ctx, &pizza.GetCartHistoryRequest{})
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, resp)
}

func (h *Handler) GetCartItemHistory(ctx *gin.Context) {

	id := ctx.GetInt("id")

	resp, err := h.GRPCClient.Pizza().GetCartItemHistory(ctx, &pizza.GetCarItemtHistoryRequest{
		CartHistoryId: int32(id),
	})
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, resp)
}
