package handler

import (
	"fmt"
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/generated/pizza"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Again, split your code into different files (orders, cart etc).
// This will make your code cleaner and easier to read.

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

// Suggestion: In GetPizzaById, DeletePizza you use strconv.Atoi without handling the error.
// If a user passes a non-numeric id, your handler will use id = 0, which may cause confusion or unwanted behavior.
// It's safer to check the error and return a Bad Request if parsing fails.
func (h *Handler) GetPizzaById(ctx *gin.Context) {

	idStr := ctx.Param("id")
	typeIdStr := ctx.Param("typeId")

	id, _ := strconv.Atoi(idStr) // <-- Handle the error here.
	typeId, _ := strconv.Atoi(typeIdStr) // <-- Handle the error here.

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
	typeIdStr := ctx.Param("typeId")

	id, _ := strconv.Atoi(idStr) // <-- Handle the error here.
	typeId, _ := strconv.Atoi(typeIdStr) // <-- Handle the error here.

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

	idStr := ctx.Param("id")
	req.Id, _ = strconv.Atoi(idStr) // <-- Handle the error here.

	resp, err := h.GRPCClient.Pizza().Cart(ctx, &pizza.CartRequest{
		Items:  req.Items,
		UserId: int32(req.UserId),
		Id:     int32(req.Id),
	})
	if err != nil {
		fmt.Println("err:", err) // <-- Avoid debug prints, use a logger if needed.
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

func (h *Handler) UpdatePizzaInCart(ctx *gin.Context) {

	var req pizza.CartItems

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(c.Empty, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().UpdatePizzaInCart(ctx, &req)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(c.OK, resp)
}

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
