package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/xiaoyi-byte/student-demo/dal"
	demo "github.com/xiaoyi-byte/student-demo/kitex_gen/demo/studentservice"
	"log"
	"net"
	_ "net/http/pprof"
)

func main() {
	//go http.ListenAndServe("localhost:8080", nil)
	dal.Init()
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"}) // r should not be reused.
	if err != nil {
		log.Fatal(err)
	}
	addr, _ := net.ResolveTCPAddr("tcp", ":9999")
	svr := demo.NewServer(new(StudentServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "student"}),
		server.WithRegistry(r),
		server.WithServiceAddr(addr))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
