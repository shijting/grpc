package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/shijting/grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net/http"
)

func main()  {
	cert,_:=tls.LoadX509KeyPair("cert/client.pem","cert/client.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds:=credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},// 加载客户端证书
		ServerName: "localhost",
		RootCAs:      certPool,
	})

	gwmux:=runtime.NewServeMux()
	opt:=[]grpc.DialOption{grpc.WithTransportCredentials(creds)}
	// 注册prod的http服务
	err:=services.RegisterProdServiceHandlerFromEndpoint(context.Background(),gwmux,"localhost:8081",opt)
	if err != nil {
		log.Fatal(err)
	}

	// 注册order的http
	err=services.RegisterOrderServiceHandlerFromEndpoint(context.Background(),gwmux,"localhost:8081",opt)
	if err != nil {
		log.Fatal(err)
	}
	httpServer := &http.Server{
		Addr: ":8080",
		Handler: gwmux,
	}
	err = httpServer.ListenAndServe()
	if err !=nil {
		log.Fatal(err)
	}
}
