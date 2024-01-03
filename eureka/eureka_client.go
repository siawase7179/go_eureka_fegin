package eureka

import (
	"fmt"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/sirupsen/logrus"
)

var client *eureka.Client
var instance *eureka.InstanceInfo

func init() {

}

func Init(eurekaURL []string) {
	client = eureka.NewClient(eurekaURL)
}

func InitFromFile(filePath string) error {
	var err error

	client, err = eureka.NewClientFromFile(filePath)
	if err != nil {
		return err
	}
	return nil
}

func NewInstance(serviceName string, hostname string, port int) error {
	var ip string = hostname

	instance = eureka.NewInstanceInfo(ip, serviceName, ip, port, 30, false)
	instance.InstanceID = fmt.Sprintf("%s:%d", ip, port)
	instance.SecurePort = &eureka.Port{Port: port, Enabled: false}
	instance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}

	err := client.RegisterInstance(serviceName, instance)
	if err != nil {
		return err
	}

	logrus.Info("eureka client init")

	return nil
}

func HeartBeat() error {
	return client.SendHeartbeat(instance.App, instance.InstanceID)
}

func GetApplication(appID string) (*eureka.Application, error) {
	app, err := client.GetApplication(appID)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func Unregister() error {
	err := client.UnregisterInstance(instance.App, instance.InstanceID)
	if err != nil {
		logrus.Error("AppId:" + instance.App + ",InstanceID:" + instance.InstanceID + ",error:" + err.Error())
	}
	logrus.Info("eureka destry")
	return err
}
