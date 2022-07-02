package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/config"
	"github.com/pavel/gateway/pkg/user"
	"github.com/pavel/gateway/pkg/workout"
	"github.com/sirupsen/logrus"
	"log"
	"path"
	"runtime"
)

func main() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Info("Init logrus logger")
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
	})
	l.Info("test")

	err, cfg := config.InitConfig(l)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("cfg: %v", cfg)
	r := gin.Default()
	r.Use(CORSMiddleware)
	api := r.Group("api")

	userSvc := user.RegisterRoute(api, *cfg)
	workout.RegisterRoute(api, *cfg, userSvc)

	r.Run(":" + cfg.Port)
}

func CORSMiddleware(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, access-control-allow-origin, access-control-allow-headers, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

	ctx.Next()
}
