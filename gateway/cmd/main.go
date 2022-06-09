package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/config"
	"github.com/pavel/gateway/pkg/user"
	"log"
)

func main() {
	err, cfg := config.InitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("cfg: %v", cfg)
	r := gin.Default()

	user.RegisterRoute(r, *cfg)

	r.Run(":" + cfg.Port)
}
