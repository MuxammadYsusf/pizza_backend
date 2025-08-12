package handler

import (
	"fmt"
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/generated/pizza"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreatePizzaType(ctx *gin.Context) {
	var req pizza.CreatePizzaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().CreatePizzaType(ctx, &req)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) CreatePizza(ctx *gin.Context) {
	var req pizza.CreatePizzaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().CreatePizza(ctx, &req)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) GetPizzas(ctx *gin.Context) {

	resp, err := h.GRPCClient.Pizza().GetPizzas(ctx, &pizza.GetPizzasRequest{})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (h *Handler) GetPizzaById(ctx *gin.Context) {

	idStr := ctx.Param("id")
	typeIdStr := ctx.Param("typeId")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}
	typeId, err := strconv.Atoi(typeIdStr)
	if err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().GetPizzaById(ctx, &pizza.GetPizzaByIdRequest{
		Id:     int32(id),
		TypeId: int32(typeId),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) UpdatePizza(ctx *gin.Context) {

	var req pizza.UpdatePizzaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(c.Empty, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().UpdatePizza(ctx, &req)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) DeletePizza(ctx *gin.Context) {

	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("HERE?")
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().DeletePizza(ctx, &pizza.DeletePizzaRequest{
		Id: int32(id),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) GetPizzaCost(ctx *gin.Context) {

	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().GetPizzaCost(ctx, &pizza.CartItems{
		PizzaId: int32(id),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}
