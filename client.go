package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/shijting/grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"time"
)

func main()  {
	/** grpc.WithInsecure()不使用https 双向认证
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	 */
	// 双向认证， 设置客户端证书
	cert,_:=tls.LoadX509KeyPair("cert/client.pem","cert/client.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds:=credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},// 加载客户端证书
		ServerName: "localhost",
		RootCAs:      certPool,
	})
	// 连接服务端
	//conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(creds))
	//defer conn.Close()
	//prodClient := services.NewProdServiceClient(conn)
	//prodRes, err :=prodClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 13})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(prodRes.ProdStock)

	conn,err:=grpc.Dial("localhost:8081",grpc.WithTransportCredentials(creds))
	if err!=nil{
		log.Fatal(err)
	}
	defer conn.Close()

	prodClient:=services.NewProdServiceClient(conn)
	ctx := context.Background()
	// 根据id获取商品库存
	prodRes,err:=prodClient.GetProdStock(ctx,&services.ProdRequest{ProdId:12, ProdArea: services.ProdAreas_B})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(prodRes.ProdStock)
	// 获取多个商品库存
	prodsRes,err := prodClient.GetProdStocks(ctx,&services.QuerySize{Size: 6})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(prodsRes.Prodres)
	// 获取商品信息
	info, err := prodClient.GetProdInfo(ctx, &services.ProdRequest{ProdId:12})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(info)

	// 订单相关
	fmt.Println()
	orderClient := services.NewOrderServiceClient(conn)
	// 生成新订单
	order := services.OrderMain{
		UserId:     99,
		OrderNo:    "20200623i2341238990",
		OrderMoney: 998.9,
		OrderTime:  &timestamp.Timestamp{Seconds: time.Now().Unix()},
	}
	newOrderRet,err := orderClient.NewOrder(ctx, &services.OrderRequest{OrderMain: &order})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(newOrderRet)
}
