package feign

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/sirupsen/logrus"
)

var mutex sync.Mutex
var feignMap = make(map[string]*feignClient)

type feignClient struct {
	App   *eureka.Application
	Index int
}

type feginError struct {
	Message string
}

func (e *feginError) Error() string {
	return e.Message
}

func Append(app eureka.Application) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, instance := range app.Instances {
		if instance.Status == "UP" {
			feignMap[app.Name] = &feignClient{App: &app, Index: 0}
		}
	}
}

func getNextInstance(feign *feignClient) (*eureka.InstanceInfo, error) {
	instances := feign.App.Instances

	for i := 0; i < len(instances); i++ {
		feign.Index = (feign.Index + 1) % len(instances)

		if instances[feign.Index].Status == "UP" {
			return &instances[feign.Index], nil
		}
	}

	return nil, errors.New("No 'UP' instances found")
}

func getInstanceInfo(appId string) (*eureka.InstanceInfo, error) {
	mutex.Lock()
	defer mutex.Unlock()
	feign := feignMap[appId]
	if feign == nil {
		return nil, &feginError{Message: "feign not found"}
	}

	instance, err := getNextInstance(feign)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

type RequeustOption struct {
	Method string
	Path   string
	Body   string
	Header map[string]string
}

type FeignResponse struct {
	Response http.Response
	Body     []byte
}

func Request(appId string, option RequeustOption) (*FeignResponse, error) {
	ins, err := getInstanceInfo(appId)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:%d%s", ins.IpAddr, ins.Port.Port, option.Path)
	logrus.Info("url:" + url)

	request, err := http.NewRequest(option.Method, url, bytes.NewBuffer([]byte(option.Body)))
	if err != nil {
		return nil, err
	}

	for key, value := range option.Header {
		request.Header.Add(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &FeignResponse{Response: *response, Body: body}, nil
}
