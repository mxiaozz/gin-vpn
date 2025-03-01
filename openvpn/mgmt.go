package openvpn

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"openvpn.funcworks.net/config"
	"openvpn.funcworks.net/domain"
	"openvpn.funcworks.net/log"
)

var Mgmt = ovpnMgmt{
	url:     config.Viper.GetString("ovpn.mgmt"),
	reqChan: make(chan *domain.MgmtRequest, 1),
	rspChan: make(chan *domain.MgmtRequest, 1),
}

type ovpnMgmt struct {
	mutex sync.Mutex

	url     string
	conn    net.Conn
	reqChan chan *domain.MgmtRequest
	rspChan chan *domain.MgmtRequest

	request *domain.MgmtRequest
}

func (m *ovpnMgmt) Run() {
	go Mgmt.write()
	go Mgmt.read()
}

func (m *ovpnMgmt) write() {
	for req := range m.reqChan {
		if m.conn == nil {
			log.Warnf("未有 mgmt 连接，丢弃 %s 指令操作", req.Command)
			continue
		}
		m.request = req

		m.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		_, err := m.conn.Write([]byte(req.Command + "\n"))
		if err != nil {
			log.Errorf(err.Error())
			m.conn.Close()
			m.conn = nil
		}
	}
}

func (m *ovpnMgmt) read() {
	err := m.checkMgmtAddr()
	if err != nil {
		return
	}

	var size = 1024
	var data bytes.Buffer
	tmp := make([]byte, size)

	for {
		Mgmt.conn, err = net.DialTimeout("tcp", Mgmt.url, 3*time.Second)
		if err != nil {
			log.Errorf("openvpn mgmt: %s", err.Error())
			Mgmt.conn = nil
			time.Sleep(30 * time.Second)
			continue
		}

		// config - 减少客户端连接的 ENV 输出（暂时不需要的信息）
		go m.Execute("env-filter 3")

		for {
			n, err := m.conn.Read(tmp)
			if err != nil {
				Mgmt.conn.Close()
				Mgmt.conn = nil
				break
			}
			data.Write(tmp[:n])

			if n < size || (tmp[size-2] == 13 && tmp[size-10] == 10) {
				// filter
				rst := m.filterAutoOutput(data)

				// clear
				data.Reset()

				// read next
				if rst == "" {
					continue
				}

				// respone to caller
				if m.request != nil {
					log.Debugf("openvpn mgmt: command = %s, respone = %s", m.request.Command, rst)
					m.request.Respone = rst
					m.rspChan <- m.request
				} else {
					log.Debugf("openvpn mgmt: command = %s, respone = %s", "", rst)
				}
			}
		}
	}
}

func (m *ovpnMgmt) checkMgmtAddr() error {
	_, err := net.ResolveTCPAddr("tcp4", Mgmt.url)
	if err != nil {
		log.Errorf("openvpn mgmt 地址配置错误")
		return err
	}
	return nil
}

func (m *ovpnMgmt) filterAutoOutput(data bytes.Buffer) string {
	rst := data.String()

	writer := bytes.NewBuffer([]byte{})
	reader := bufio.NewReader(strings.NewReader(rst))
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		if line[0] != 62 { // == ">"
			writer.Write(line)
			writer.WriteByte(13)
			writer.WriteByte(10)
		}
	}
	rst = writer.String()

	return rst
}

func (m *ovpnMgmt) Execute(cmd string) (*domain.MgmtRequest, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	req := &domain.MgmtRequest{Command: cmd}
	m.reqChan <- req

	select {
	case r := <-m.rspChan:
		return r, nil
	case <-time.After(10 * time.Second):
		return nil, errors.New("read timeout")
	}
}
