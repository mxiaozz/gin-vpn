package openvpn

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"openvpn.funcworks.net/config"
	"openvpn.funcworks.net/log"
	rsp "openvpn.funcworks.net/respone"
)

var PKI = &OpenVpnPKI{
	CaPath:  config.Viper.GetString("ovpn.pki"),
	RootDir: filepath.Dir(config.Viper.GetString("ovpn.config")),
	PKIDir:  filepath.Dir(config.Viper.GetString("ovpn.pki")),
}

type OpenVpnPKI struct {
	CaPath  string
	RootDir string
	PKIDir  string
}

func (k *OpenVpnPKI) Init(ctx *gin.Context) {
	b, err := k.isInited()
	if err != nil {
		log.Errorf("查找ca.crt文件失败: %s", err.Error())
		rsp.Fail("PKI初始化失败", ctx)
		return
	}
	if b {
		rsp.Fail("PKI已初始化过，不能再次执行", ctx)
		return
	}

	cmd := exec.Command("ovpn_initpki", "nopass")
	cmd.Dir = k.RootDir
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Errorf("PKI初始化失败: %s", err.Error())
		rsp.Fail("PKI初始化失败", ctx)
		return
	}

	rsp.Ok(ctx)
}

func (k *OpenVpnPKI) Reset(ctx *gin.Context) {
	b, err := k.isInited()
	if err != nil {
		log.Errorf("查找ca.crt文件失败: %s", err.Error())
		rsp.Fail("重置PKI失败", ctx)
		return
	}

	if b {
		cmd := exec.Command("rm", "-rf", k.PKIDir)
		log.Debugf("重置PKI: %s", cmd.String())

		err := cmd.Run()
		if err != nil {
			log.Errorf(err.Error())
			rsp.Fail("重置PKI失败", ctx)
			return
		}
	}

	k.Init(ctx)
}

func (k *OpenVpnPKI) isInited() (bool, error) {
	_, err := os.Stat(k.CaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		log.Debugf(err.Error())
		return false, err
	}

	return true, nil
}

func (k *OpenVpnPKI) Status(ctx *gin.Context) {
	b, err := k.isInited()
	if err != nil {
		rsp.Fail("PKI状态查询失败", ctx)
		return
	}
	rsp.OkWithData(b, ctx)
}
