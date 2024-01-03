package test

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/siawase7179/go_eureka_fegin/eureka"
	feign "github.com/siawase7179/go_eureka_fegin/eureka/fegin"
)

const (
	serviceName = "GO-SERVICE"
)

func init() {
	eureka.Init([]string{"http://localhost:8761/eureka"})

	eureka.NewInstance(serviceName, "localhost", 8082)
}

func TestEurekaHeartBeat(t *testing.T) {
	err := eureka.HeartBeat()
	if err != nil {
		t.Fatal(err)
	}
}

func TestEurekaUnRegistry(t *testing.T) {
	err := eureka.Unregister()
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetInstance(t *testing.T) {
	app, err := eureka.GetApplication(serviceName)
	if err != nil {
		t.Fatal(err)
	}

	json, err := json.Marshal(app)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(json))
}

type Response struct {
	Message string `json:"message"`
}

func TestFegin(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer listener.Close()

	go func() {
		err := http.Serve(listener, handler)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			t.Error(err)
			return
		}
	}()

	time.Sleep(1 * time.Second)

	app, err := eureka.GetApplication(serviceName)
	if err != nil {
		t.Fatal(err)
	}
	feign.Append(*app)

	response, err := feign.Request(serviceName, feign.RequeustOption{
		Method: "GET",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(response.Response.Status)
}
