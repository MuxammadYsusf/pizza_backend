package handler

import (
	"fmt"
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/generated/pizza"

	"github.com/gin-gonic/gin"
)

func (h *Handler) OrderPizza(ctx *gin.Context) {

	var req struct {
		Items  []*pizza.OrderItemData `json:"items"`
		UserId int                    `json:"userId"`
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

	fmt.Println("req", req)

	resp, err := h.GRPCClient.Order().OrderPizza(ctx, &pizza.OrderPizzaRequest{
		UserId: int32(req.UserId),
		Items:  req.Items,
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}
