package main

import (
	"fmt"
	"github.com/huiming23344/balanceapi/config"
	"github.com/huiming23344/balanceapi/routers"
	"log"
	"net/http"
)

func main() {
	router := routers.InitRouter()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic("load config fail: " + err.Error())
	}
	config.SetGlobalConfig(cfg)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Printf("Listen: %s\n", err)
	}

}
