package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/shijting/grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
)

func main()  {
	// tls 是ssl的升级版
	cert,_:=tls.LoadX509KeyPair("cert/server.pem","cert/server.key")
	// 证书池
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds:=credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, // 服务端证书
		ClientAuth:   tls.RequireAndVerifyClientCert, // 验证客户端证书(双向认证)
		ClientCAs:    certPool,
	})


	grpcServer := grpc.NewServer(grpc.Creds(creds))
	// 注册商品服务
	services.RegisterProdServiceServer(grpcServer, new(services.ProdService))
	// 订单服务
	services.RegisterOrderServiceServer(grpcServer, new(services.OrderService))
	// 其他服务
	lis,_ :=net.Listen("tcp", ":8081")
	grpcServer.Serve(lis)

	/** 让grpc提供http服务
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		grpcServer.ServeHTTP(writer, request)
	})
	httpServer := &http.Server{
		Addr: ":8081",
		Handler: mux,
	}

	httpServer.ListenAndServe()

	 */
}
