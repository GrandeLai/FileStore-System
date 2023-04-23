package main

import (
	"flag"
	"fmt"

	"FileStore-System/service/share/api/internal/config"
	"FileStore-System/service/share/api/internal/handler"
	"FileStore-System/service/share/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "D:\\golangPro\\FileStore-System\\service\\share\\api\\etc\\share-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
