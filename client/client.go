package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/whalethinker/server_with_register/env"
	"github.com/whalethinker/server_with_register/http"
	"log"
	"math/rand/v2"
	http2 "net/http"
	"sync"
	"time"
)

type Client interface {
	Call(path string, method string, headers map[string]string, params map[string]string, body string) ([]byte, error)
}

type ClientImpl struct {
	Psm             string
	PodInfoList     []*PodInfo
	LastRefreshTime time.Time
	mtx             sync.Mutex
}

func (c *ClientImpl) Call(path string, method string, headers map[string]string, params map[string]string, body string) ([]byte, error) {
	if time.Now().After(c.LastRefreshTime.Add(10 * time.Second)) {
		go func() {
			err := c.refresh()
			if err != nil {
				log.Println("refresh err:", err)
			}
		}()
	}
	podInfo := c.PodInfoList[int(RandomInt64())%len(c.PodInfoList)]
	realPath := fmt.Sprintf("http://%s%s", podInfo.Addr, path)
	return http.Call(realPath, method, headers, params, body)
}

func RandomInt64() int64 {
	return rand.Int64()
}

func (c *ClientImpl) refresh() error {
	podInfoList, err := GetPodMap(c.Psm)
	if err != nil {
		return err
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.PodInfoList = podInfoList
	c.LastRefreshTime = time.Now()
	return nil
}

func NewClient(ctx context.Context, psm string) (Client, error) {
	podInfoList, err := GetPodMap(psm)
	if err != nil {
		return nil, err
	}
	return &ClientImpl{
		Psm:             psm,
		PodInfoList:     podInfoList,
		LastRefreshTime: time.Now(),
		mtx:             sync.Mutex{},
	}, nil
}

func GetPodMap(psm string) ([]*PodInfo, error) {
	dsAddr := env.DSAddr()
	url := fmt.Sprintf("http://%v/get_psm_pod_list", dsAddr)
	params := map[string]string{
		"psm": psm,
	}
	resp, err := http.Call(url, http2.MethodGet, make(map[string]string), params, "")
	if err != nil {
		return []*PodInfo{}, errors.Wrap(err, "Register failed")
	}
	respMap := map[string][]*PodInfo{}
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return []*PodInfo{}, errors.Wrap(err, "Unmarshal failed")
	}
	podList, ok := respMap["PodList"]
	if !ok {
		return []*PodInfo{}, errors.New("PodList not found")
	}
	return podList, nil
}

type PodInfo struct {
	Addr string
	PSM  string
}
