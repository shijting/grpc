package services

import (
	"context"
	"fmt"
)

type OrderService struct {}

func (this *OrderService) NewOrder(ctx context.Context, request *OrderRequest) (*OrderResponse, error) {
	fmt.Println(request.OrderMain.UserId)
	ret := &OrderResponse{
		Status:        "ok",
		Msg:           "创建新订单成功",
	}
	return ret, nil
}
