https://grpc.io/docs/languages/go/quickstart/

### gRPC
#### 安装文档
[文档](https://grpc.io/docs/languages/go/quickstart/)
#### 客户端调试工具
[evans](https://github.com/ktr0731/evans)

#### 示例
```shell
evans --host localhost --port 9090 -r repl
show service # 查看服务
call CreateUser # 调用服务
```

#### protoc 
[文档和下载地址](https://github.com/protocolbuffers/protobuf)

1. 下载解压后将bin目录加入环境变量中或者复制bin目录下的protoc 文件到GOPATH(go env GOPATH)的bin目录中(需要将GOPATH的bin目录添加到环境变量)
2. 将include\google目录复制到你项目的放proto/google文件的目录中 ，比如本项目的proto/google。 或者将 "include "目录的内容也复制到某个地方，例如
复制到"/usr/local/include/"中, 可以将 include 放到任意位置，使用protoc 生成中间文件的时候 -I 进去即可， 如：`protoc -I /usr/local/include`


#### grpc-gateway
[grpc-gateway github](https://github.com/grpc-ecosystem/grpc-gateway)

[grpc-gateway文档](https://grpc-ecosystem.github.io/grpc-gateway/)

使用示例：

### 安装
新建tools/tools.go
```go
//go:build tools
// +build tools
package tools

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
```

或者手动go get 安装

因为我们不直接在代码中使用它们，而只是想将他们安装到本地机器上
然后执行
```go
go mod tidy

go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
This will place four binaries in your $GOBIN;

##### 下载并复制下面proto文件放到proto目录下
[googleapis](https://github.com/googleapis/googleapis)
```shell
proto/google/api/annotations.proto
proto/google/api/field_behavior.proto
proto/google/api/http.proto
proto/google/api/httpbody.proto
```

在rpc服务(service_simple_bank.proto)中添加注解
```go
syntax = "proto3";

package pb;

// 无需导入user.proto，因为在rpc_create_user.proto中已经导入了
//import "user.proto";
import "rpc_create_user.proto";
import "rpc_login_user.proto";
import "google/api/annotations.proto";

option go_package = "github.com/shijting/learngo/pb";

service SimpleBank {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse){
    option (google.api.http) = {
      post: "/v1/create_user"
      body: "*"
    };
  };
  rpc LoginUser(LoginUserRequest) returns (LoginUserResponse){
    option (google.api.http) = {
      post: "/v1/login_user"
      body: "*"
    };
  };
}
```

##### 生成pb文件
```go
protoc -I proto --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
        proto/*.proto
```

##### 生成swagger文件
[在官方项目中下载](https://github.com/grpc-ecosystem/grpc-gateway/tree/main/protoc-gen-openapiv2)

把protoc-gen-openapiv2/options/annotations.proto和protoc-gen-openapiv2/options/openapiv2.proto复制到proto目录下
```shell
proto/protoc-gen-openapiv2/options/annotations.proto
proto/protoc-gen-openapiv2/options/openapiv2.proto
```
[proto文件示例](https://github.com/grpc-ecosystem/grpc-gateway/blob/main/examples/internal/proto/examplepb/a_bit_of_everything.proto)


在service_simple_bank.proto中添加注解

```go

import "protoc-gen-openapiv2/options/annotations.proto";

option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  info: {
    title: "Simple Bank API";
    version: "1.0";
    contact: {
      name: "gRPC-Gateway project";
      url: "https://github.com/grpc-ecosystem/grpc-gateway";
      email: "none@example.com";
    };
    license: {
      name: "BSD 3-Clause License";
      url: "https://github.com/grpc-ecosystem/grpc-gateway/blob/main/LICENSE.txt";
    };
    extensions: {
      key: "x-something-something";
      value {string_value: "yadda"}
    }
  };
};
```
添加注解后，添加openapiv2插件，重新生成pb文件

```go
--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank
```
完整的
```go
protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
proto/*.proto
```

##### 生成swagger-ui
[下载swagger-ui](https://github.com/swagger-api/swagger-ui)

```shell
git clone https://github.com/swagger-api/swagger-ui.git
```
- 把swagger-ui/dist目录下的文件复制到doc/swagger目录下
- 修改swagger-initializer.js文件的url为openapiv2生成的json文件

```shell
url: "simple_bank.swagger.json"
```

把swagger ui嵌入到我们的网络服务器中，在main.go的runGatewayServer中添加添加文件服务
```go
fs := http.FileServer(http.Dir("./doc/swagger"))
mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
```

在浏览器中访问：http://localhost:8080/swagger/

###### 优化 - 把静态前端文件打包成二进制文件， 把swagger-ui嵌入到二进制文件中

我们在部署的时候必须把swagger-ui的文件也一起部署，这样不太方便，我们可以把swagger-ui的文件嵌入到二进制文件中，这样部署的时候就不需要再部署swagger-ui的文件了
这个二进制文件会被加载到服务器的内存中

[statik地址](https://github.com/rakyll/statik)
```go
go get github.com/rakyll/statik
go install github.com/rakyll/statik
statik --hlep
statik -src=doc/swagger -dest=doc
```

修改main.go的runGatewayServer
```go
import (
    _ "github.com/shijting/learngo/doc/statik"
)


statikFS, err := fs.New()
if err != nil {
    log.Fatal().Err(err).Msg("cannot create statik fs")
}

swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
mux.Handle("/swagger/", swaggerHandler)
```
