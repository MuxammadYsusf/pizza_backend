package handler

import (
	pb "github/http/copy/task4/genproto"
	"github/http/copy/task4/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreatePizzaType(ctx *gin.Context) {
	var req pb.CreatePizzaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().CreatePizzaType(ctx, &req)
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) CreatePizza(ctx *gin.Context) {
	var req pb.CreatePizzaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().CreatePizza(ctx, &req)
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) GetPizzas(ctx *gin.Context) {

	resp, err := h.GRPCClient.Pizza().GetPizzas(ctx, &pb.GetPizzasRequest{})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) GetPizzaById(ctx *gin.Context) {

	idStr := ctx.Param("id")
	typeIdStr := ctx.Param("typeId")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}
	typeId, err := strconv.Atoi(typeIdStr)
	if err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().GetPizzaById(ctx, &pb.GetPizzaByIdRequest{
		Id:     int32(id),
		TypeId: int32(typeId),
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) UpdatePizza(ctx *gin.Context) {

	var req pb.UpdatePizzaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().UpdatePizza(ctx, &req)
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) DeletePizza(ctx *gin.Context) {

	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().DeletePizza(ctx, &pb.DeletePizzaRequest{
		Id: int32(id),
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}

func (h *Handler) GetPizzaCost(ctx *gin.Context) {

	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(util.HTTPBadReq, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.Pizza().GetPizzaCost(ctx, &pb.CartItems{
		PizzaId: int32(id),
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}
