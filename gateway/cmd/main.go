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
	api := r.Group("api")

	userSvc := user.RegisterRoute(api, *cfg)
	workout.RegisterRoute(api, *cfg, userSvc)

	r.Run(":" + cfg.Port)
}
