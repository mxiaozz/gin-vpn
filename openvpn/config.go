package openvpn

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"openvpn.funcworks.net/config"
	"openvpn.funcworks.net/domain"
	"openvpn.funcworks.net/log"
	rsp "openvpn.funcworks.net/respone"
)

var Config = &OpenVpnConfig{
	Path:    config.Viper.GetString("ovpn.config"),
	RootDir: filepath.Dir(config.Viper.GetString("ovpn.config")),
}

type OpenVpnConfig struct {
	Path    string
	RootDir string
}

func (c *OpenVpnConfig) isExist() (bool, error) {
	_, err := os.Stat(c.Path)
	if err != nil {
		log.Errorf(err.Error() + "," + c.Path)
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *OpenVpnConfig) Read(ctx *gin.Context) {
	exist, err := c.isExist()
	if err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}
	if !exist {
		rsp.FailWithCode(rsp.NO_GENERATE_CONFIG, "OpenVPN服务尚未初始化配置", ctx)
		return
	}

	data, err := os.ReadFile(c.Path)
	if err != nil {
		log.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	rsp.OkWithData(string(data), ctx)
}

func (c *OpenVpnConfig) Generate(ctx *gin.Context) {
	var cfg domain.DefineConfig
	if err := ctx.ShouldBind(&cfg); err != nil {
		log.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 先生成配置
	log.Infof("开始生成OpenVPN服务配置")
	cmd := exec.Command("ovpn_genconfig")
	cmd.Dir = c.RootDir
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Args = []string{"ovpn_genconfig"}
	cmd.Args = append(cmd.Args, cfg.Params...)

	log.Debugf("开始生成配置: %s", cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	rsp.Ok(ctx)
}

func (c *OpenVpnConfig) Save(ctx *gin.Context) {
	var cfg domain.DefineConfig
	if err := ctx.ShouldBind(&cfg); err != nil {
		log.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	file, err := os.OpenFile(c.Path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	_, err = file.WriteString(cfg.Content)
	if err != nil {
		log.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	rsp.Ok(ctx)
}
