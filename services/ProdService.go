package services

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

type ProdService struct {
	
}

func (this *ProdService)GetProdStock(ctc context.Context, in *ProdRequest) (*ProdResponse, error)  {
	var stock int32 = 0
	if in.ProdArea == ProdAreas_A {
		stock = 18
	}
	if in.ProdArea == ProdAreas_B {
		stock = 37
	}
	if in.ProdArea == ProdAreas_C {
		stock = 16
	}
	fmt.Println("request", in.ProdId)
	return &ProdResponse{ProdStock: stock},nil
}
func (this *ProdService) GetProdStocks(ctx context.Context, in *QuerySize) (*ProdResponseList, error) {
	prodres := []*ProdResponse{
		&ProdResponse{ProdStock: 10},
		&ProdResponse{ProdStock: 11},
		&ProdResponse{ProdStock: 12},
		&ProdResponse{ProdStock: 13},
		&ProdResponse{ProdStock: 14},
		&ProdResponse{ProdStock: 15},
		&ProdResponse{ProdStock: 16},
	}
	return &ProdResponseList{Prodres: prodres}, nil
}
func (this *ProdService) GetProdInfo(context.Context, *ProdRequest) (*ProdModel, error) {
	ret := &ProdModel{
		ProdId: 100,
		ProdName: "测试商品",
		ProdPrice: 19.9,
		CreatedTime: &timestamp.Timestamp{Seconds: time.Now().Unix()},
	}
	return ret, nil
}