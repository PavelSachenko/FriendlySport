package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/config"
	"github.com/pavel/gateway/pkg/user"
	"github.com/pavel/gateway/pkg/workout"
	"log"
)

func main() {
	err, cfg := config.InitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("cfg: %v", cfg)
	r := gin.Default()

	userSvc := user.RegisterRoute(r, *cfg)
	workout.RegisterRoute(r, *cfg, userSvc)

	r.Run(":" + cfg.Port)
}
