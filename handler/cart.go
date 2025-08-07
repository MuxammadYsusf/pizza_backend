package handler

import (
	"fmt"
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

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(c.Empty, gin.H{"error": err.Error()})
		return
	}

	var id int
	req.Id = id

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	resp, err := h.GRPCClient.Pizza().Cart(ctx, &pizza.CartRequest{
		Items:  req.Items,
		UserId: int32(req.UserId),
		Id:     int32(id),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) IncreasePizzaInCart(ctx *gin.Context) {

	idStr := ctx.Param("id")
	pizzaIdStr := ctx.Param("typeId")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}
	pizzaId, err := strconv.Atoi(pizzaIdStr)
	if err != nil {
		ctx.JSON(c.BadReq, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("LOOK AT HERE!!!")

	resp, err := h.GRPCClient.Pizza().IncreaseAmountOfPizza(ctx, &pizza.CartItems{
		PizzaId:     int32(pizzaId),
		PizzaTypeId: int32(id),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) DecreasePizzaInCart(ctx *gin.Context) {

	idStr := ctx.Param("id")
	typeIdStr := ctx.Param("typeId")

	id, _ := strconv.Atoi(idStr)
	typeId, _ := strconv.Atoi(typeIdStr)

	fmt.Printf("id: %d \ntypeId: %d\n", id, typeId)

	resp, err := h.GRPCClient.Pizza().DeletePizza(ctx, &pizza.DeletePizzaRequest{
		Id:     int32(id),
		TypeId: int32(typeId),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) GetCartHistory(ctx *gin.Context) {

	resp, err := h.GRPCClient.Pizza().GetUserHistory(ctx, &pizza.GetCartHistoryRequest{})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) GetCartItemHistory(ctx *gin.Context) {

	id := ctx.GetInt("id")

	resp, err := h.GRPCClient.Pizza().GetCartItemHistory(ctx, &pizza.GetCarItemtHistoryRequest{
		CartHistoryId: int32(id),
	})
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err})
		return
	}

	ctx.JSON(c.OK, resp)
}
