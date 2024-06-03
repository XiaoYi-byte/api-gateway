// Code generated by hertz generator.

package main

import (
	"bytes"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/xiaoyi-byte/student-demo/hertz-gateway/biz/handler"
	"log"
	"net/http"
	"os"
	"time"
)

// customizeRegister registers customize routers.
func customizedRegister(r *server.Hertz) {
	r.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, "hertz-gateway is running")
	})

	registerGateway(r)
}

func registerGateway(r *server.Hertz) {
	group := r.Group("/gateway")
	idlPath := "./idl/student.thrift"
	content, err := os.ReadFile(idlPath)
	if err != nil {
		hlog.Fatalf("read dir error: %v", err)
	}
	svcName := "student"
	p, err := generic.NewThriftContentProvider(string(content), map[string]string{})
	if err != nil {
		hlog.Fatalf("new thrift content provider failed: %v", err)
	}
	g, err := generic.HTTPThriftGeneric(p)
	if err != nil {
		hlog.Fatalf("new http thrift generic failed: %v", err)
	}
	resolver, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}
	cli, err := genericclient.NewClient(svcName, g, client.WithResolver(resolver))
	if err != nil {
		hlog.Fatalf("new http generic client failed: %v", err)
	}
	handler.SvcMap.Store(svcName, cli)
	go func() {
		for {
			newContent, err := os.ReadFile(idlPath)
			if err != nil {
				hlog.Fatalf("read dir error: %v", err)
			}
			if !bytes.Equal(newContent, content) {
				err = p.UpdateIDL(string(newContent), map[string]string{})
				if err != nil {
					hlog.Fatalf("update IDL failed: %v", err)
				}
				hlog.Info("update IDL successfully!")
				g, err := generic.HTTPThriftGeneric(p)
				if err != nil {
					hlog.Fatalf("new http thrift generic failed: %v", err)
				}
				cli, err := genericclient.NewClient(svcName, g, client.WithResolver(resolver))
				if err != nil {
					hlog.Fatalf("new http generic client failed: %v", err)
				}
				handler.SvcMap.Store(svcName, cli)
			}
			time.Sleep(3 * time.Second)
		}
	}()

	group.POST("/:svc", handler.HttpGateway)

}

//func registerGateway(r *server.Hertz) {
//	r.POST("/gateway/:svc", handler.HttpGateway)
//}
