package router

import (
	"github.com/gin-gonic/gin"
	"openvpn.funcworks.net/openvpn"
)

func Init(engine *gin.Engine) {

	// openvpn config
	engine.GET("/getConfig", openvpn.Config.Read)
	engine.POST("/genConfig", openvpn.Config.Generate)
	engine.POST("/saveConfig", openvpn.Config.Save)

	// openvpn pki
	engine.GET("/pkiStatus", openvpn.PKI.Status)
	engine.POST("/initPKI", openvpn.PKI.Init)
	engine.POST("/resetPKI", openvpn.PKI.Reset)

	// openvpn server
	engine.POST("/serverStart", openvpn.Server.Start)
	engine.POST("/serverStop", openvpn.Server.Stop)
	engine.GET("/serverStatus", openvpn.Server.Status)
	engine.GET("/state", openvpn.Server.MgmtState)
	engine.GET("/status", openvpn.Server.MgmtStatus)

	// openvpn client
	engine.GET("/getClientCert", openvpn.Client.GetCert)
	engine.POST("/genClientCert", openvpn.Client.GenerateCert)
	engine.POST("/revokeClientCert", openvpn.Client.RevokeCert)
	engine.POST("/killClient", openvpn.Client.MgmtKill)

}
