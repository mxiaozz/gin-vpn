package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"openvpn.funcworks.net/config"
	_ "openvpn.funcworks.net/log"
	"openvpn.funcworks.net/openvpn"
	"openvpn.funcworks.net/router"
)

func main() {
	// openvpn
	if err := openvpn.Server.Run(); err == nil {
		// wait setup
		time.Sleep(5 * time.Second)
	}

	// mgmt
	openvpn.Mgmt.Run()

	// gin
	engine := gin.Default()
	router.Init(engine)

	addr := config.Viper.GetString("server.port")
	if config.Viper.GetBool("server.dev") {
		addr = "127.0.0.1:" + addr
	} else {
		addr = ":" + addr
	}

	engine.SetTrustedProxies([]string{"127.0.0.1"})
	engine.Run(addr)
}
