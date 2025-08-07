package handler

import (
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/generated/pizza"

	"github.com/gin-gonic/gin"
)

func (h *Handler) OrderPizza(ctx *gin.Context) {

	var req pizza.OrderPizzaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(c.Empty, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().OrderPizza(ctx, &req)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}
