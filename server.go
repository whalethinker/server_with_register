package server_with_register

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/pkg/errors"
	"github.com/whalethinker/logs"
	"github.com/whalethinker/server_with_register/env"
	"github.com/whalethinker/server_with_register/http"
	"log"
	http2 "net/http"
	"sync"
	"time"
)

type HertzServerWithRegister struct {
	*server.Hertz
}

func (h *HertzServerWithRegister) Spin() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		h.Hertz.Spin()
	}()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			if h.IsRunning() {
				err := register()
				if err != nil {
					logs.Error("register failed", err)
					continue
				}
			}
		}
	}()
	wg.Wait()
}

func BuildHertzServerWithRegister() (*HertzServerWithRegister, error) {
	h, err := BuildHertzServerWithCheckApi()
	if err != nil {
		return nil, err
	}
	return &HertzServerWithRegister{
		Hertz: h,
	}, nil
}

func BuildHertzServerWithCheckApi() (*server.Hertz, error) {
	// run Hertz with the consul register
	port, err := findAvailablePort()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	addr := fmt.Sprintf("%v:%v", env.PSMIP(), port)
	httpStr := fmt.Sprintf("http://%v/consul_check_ping", addr)
	registerInfo := &RegisterInfo{
		Addr:         addr,
		PSM:          env.PSM(),
		HttpCheckUrl: httpStr,
	}
	ServiceRegisterInfo = registerInfo
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(":%v", port)),
	)

	h.GET("/consul_check_ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"ping": "pong1"})
	})
	return h, nil
}

func CheckRegisterInfo() error {
	dsAddr := env.DSAddr()
	url := fmt.Sprintf("http://%v/check_pod", dsAddr)
	podInfo := &RegisterInfo{
		Addr:         ServiceRegisterInfo.Addr,
		PSM:          ServiceRegisterInfo.PSM,
		HttpCheckUrl: ServiceRegisterInfo.HttpCheckUrl,
	}
	_, err := http.Call(url, http2.MethodPost, make(map[string]string), make(map[string]string), JsonMarshal2String(podInfo))
	if err != nil {
		return errors.Wrap(err, "Register failed")
	}
	return nil
}

var ServiceRegisterInfo *RegisterInfo

type ServiceRegister struct {
}

func (s *ServiceRegister) Register(info *registry.Info) error {
	return register()
}

func register() error {
	dsAddr := env.DSAddr()
	url := fmt.Sprintf("http://%v/register", dsAddr)
	podInfo := ServiceRegisterInfo
	_, err := http.Call(url, http2.MethodPost, make(map[string]string), make(map[string]string), JsonMarshal2String(podInfo))
	if err != nil {
		return errors.Wrap(err, "Register failed")
	}
	return nil
}

func JsonMarshal2String(val any) string {
	bytes, err := json.Marshal(val)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (s *ServiceRegister) Deregister(info *registry.Info) error {
	return nil
}

type RegisterInfo struct {
	Addr         string
	PSM          string
	HttpCheckUrl string
}
