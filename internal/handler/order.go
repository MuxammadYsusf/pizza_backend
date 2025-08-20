package handler

import (
	pb "github/http/copy/task4/genproto"
	"github/http/copy/task4/pkg/util"

	"github.com/gin-gonic/gin"
)

func (h *Handler) OrderPizza(ctx *gin.Context) {

	var req struct {
		Items  []*pb.OrderItemData `json:"items"`
		Limit  int32               `json:"limit"`
		UserId int                 `json:"userId"`
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

	resp, err := h.GRPCClient.Order().OrderPizza(ctx, &pb.OrderPizzaRequest{
		Items:  req.Items,
		UserId: int32(req.UserId),
		Limit:  req.Limit,
	})
	if err != nil {
		ctx.JSON(util.HTTPServerErr, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(util.HTTPOK, resp)
}
