package openvpn

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/process"
	"openvpn.funcworks.net/config"
	"openvpn.funcworks.net/log"
	rsp "openvpn.funcworks.net/respone"
)

var Server = &OpenVpnServer{}

type OpenVpnServer struct {
}

func (s *OpenVpnServer) Run() error {
	b, err := Config.isExist()
	if err != nil {
		return err
	}
	if !b {
		return errors.New("OpenVPN服务尚未初始化配置")
	}

	cmd := exec.Command("ovpn_run")
	cmd.Dir = filepath.Dir(config.Viper.GetString("ovpn.config"))
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debugf("启动OpenVPN服务: %s", cmd.String())
	err = cmd.Start()
	if err != nil {
		return err
	}

	// 防止僵尸进程
	go cmd.Wait()

	return nil
}

func (s *OpenVpnServer) Start(ctx *gin.Context) {
	err := s.Run()
	if err != nil {
		log.Errorf(err.Error())
		rsp.Fail("启动失败", ctx)
		return
	}

	rsp.Ok(ctx)
}

func (s *OpenVpnServer) Stop(ctx *gin.Context) {
	processes, err := process.Processes()
	if err != nil {
		log.Errorf("查找系统进程失败: %s", err.Error())
		rsp.Fail("查找系统进程失败", ctx)
		return
	}

	var pid string
	for _, p := range processes {
		userName, _ := p.Username()
		if userName == "nobody" {
			pid = strconv.Itoa(int(p.Pid))
			break
		}
	}
	log.Debugf("openvpn pid: %s", pid)

	if pid == "" {
		rsp.Fail("服务未启动", ctx)
		return
	}

	cmd := exec.Command("kill", pid)
	cmd.Env = os.Environ()
	log.Debugf("停止服务: %s", cmd.String())
	out, err := cmd.CombinedOutput()
	log.Debugf(string(out))
	if err != nil {
		log.Errorf("服务停止失败: %s", err.Error())
		rsp.Fail("服务停止失败", ctx)
		return
	}

	rsp.Ok(ctx)
}

func (s *OpenVpnServer) Status(ctx *gin.Context) {
	processes, err := process.Processes()
	if err != nil {
		log.Errorf("查找系统进程失败: %s", err.Error())
		rsp.OkWithData("unknown", ctx)
		return
	}

	var pid string
	for _, p := range processes {
		userName, _ := p.Username()
		if userName == "nobody" {
			pid = strconv.Itoa(int(p.Pid))
			break
		}
	}
	log.Debugf("openvpn pid: %s", pid)

	if pid == "" {
		b, err := Config.isExist()
		if err != nil {
			log.Errorf(err.Error())
			rsp.OkWithData("unknown", ctx)
			return
		}
		if !b {
			log.Errorf("OpenVPN服务尚未初始化配置")
			rsp.OkWithData("notinit", ctx)
			return
		}

		rsp.OkWithData("stopped", ctx)
	} else {
		rsp.OkWithData("running", ctx)
	}
}

func (s *OpenVpnServer) MgmtState(ctx *gin.Context) {
	obj, err := Mgmt.Execute("state")
	if err != nil {
		log.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	rsp.OkWithData(obj, ctx)
}

func (s *OpenVpnServer) MgmtStatus(ctx *gin.Context) {
	obj, err := Mgmt.Execute("status")
	if err != nil {
		log.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	rsp.OkWithData(obj, ctx)
}
