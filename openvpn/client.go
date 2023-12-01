package openvpn

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openvpn.funcworks.net/config"
	"openvpn.funcworks.net/domain"
	"openvpn.funcworks.net/log"
	rsp "openvpn.funcworks.net/respone"
)

var Client = &OpenVpnClient{
	RootDir: filepath.Join(filepath.Dir(config.Viper.GetString("ovpn.pki")), "issued"),
}

type OpenVpnClient struct {
	RootDir string
}

func (c *OpenVpnClient) GetCert(ctx *gin.Context) {
	userName := ctx.Query("name")
	if userName == "" {
		rsp.Fail("name不能为空", ctx)
		return
	}
	userName = strings.Split(userName, "/")[0]

	certPath := filepath.Join(c.RootDir, userName+".crt")
	_, err := os.Stat(certPath)
	if err != nil {
		log.Errorf(err.Error() + "," + certPath)

		if os.IsNotExist(err) {
			rsp.FailWithCode(rsp.NO_CLIENT_CERT, "证书不存在", ctx)
			return
		}

		rsp.Fail(err.Error(), ctx)
		return
	}

	cmd := exec.Command("ovpn_getclient", userName)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Debugf(string(out))
		log.Errorf(err.Error())
		rsp.Fail("读取证书失败", ctx)
		return
	}

	client := domain.Client{
		UserName:    userName,
		CertContent: string(out),
	}

	rsp.OkWithData(client, ctx)
}

func (c *OpenVpnClient) GenerateCert(ctx *gin.Context) {
	client, err := c.getRequestClient(ctx)
	if err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}
	if client.CertExpire <= 0 {
		client.CertExpire = 365
	}

	cmd := exec.Command("easyrsa", "build-client-full", client.UserName, "nopass")
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader("yes")

	// 证书有效期
	cmd.Env = append(cmd.Env, "EASYRSA_CERT_EXPIRE="+strconv.Itoa(client.CertExpire))

	log.Debugf(cmd.String())
	err = cmd.Run()
	if err != nil {
		log.Errorf(err.Error())
		rsp.Fail("生成用户证书失败", ctx)
		return
	}

	rsp.Ok(ctx)
}

func (c *OpenVpnClient) RevokeCert(ctx *gin.Context) {
	client, err := c.getRequestClient(ctx)
	if err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}

	cmd := exec.Command("ovpn_revokeclient", client.UserName)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader("yes")
	err = cmd.Run()
	if err != nil {
		log.Errorf(err.Error())
		rsp.Fail("用户证书吊销失败", ctx)
		return
	}

	rsp.Ok(ctx)
}

func (c *OpenVpnClient) MgmtKill(ctx *gin.Context) {
	client, err := c.getRequestClient(ctx)
	if err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}

	obj, err := Mgmt.Execute("kill " + client.UserName)
	if err != nil {
		log.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	if strings.HasPrefix(obj.Respone, "ERROR:") {
		rsp.Fail(obj.Respone, ctx)
		return
	}

	rsp.OkWithData(obj, ctx)
}

func (c *OpenVpnClient) getRequestClient(ctx *gin.Context) (domain.Client, error) {
	var client domain.Client
	if err := ctx.ShouldBind(&client); err != nil {
		log.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return client, err
	}
	if client.UserName == "" {
		return client, errors.New("name不能为空")
	}
	client.UserName = strings.Split(client.UserName, "/")[0]

	return client, nil
}
