package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	"log"
	"net"
	"net/http"
	"os"
)

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("failed to find an external IP address")
}

func registerInConsul(consulAddr string) (string, error) {
	cfg := capi.DefaultConfig()
	cfg.Address = consulAddr
	client, err := capi.NewClient(cfg)
	if err != nil {
		return "", err
	}
	host, err := externalIP()
	if err != nil {
		return "", err
	}
	serviceID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	err = client.Agent().ServiceRegister(&capi.AgentServiceRegistration{
		ID:      serviceID.String(),
		Name:    "messaging-service",
		Port:    8080,
		Address: host,
		Check: &capi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:8080/healthcheck", host),
			Interval: "10s",
			Timeout:  "30s",
		},
	})
	if err != nil {
		return "", err
	}
	//keyValuePairs := client.KV()
	//kafkaAddrKV, _, err := keyValuePairs.Get("kafka_address", &capi.QueryOptions{})
	//if err != nil {
	//	return err
	//}
	//kafkaTopicKV, _, err := keyValuePairs.Get("kafka_topic", &capi.QueryOptions{})
	//if err != nil {
	//	return err
	//}
	//kafkaService = string(kafkaAddrKV.Value)
	//topic = string(kafkaTopicKV.Value)
	return serviceID.String(), nil
}

func unregisterConsul(consulAddr string, serviceID string) {
	cfg := capi.DefaultConfig()
	cfg.Address = consulAddr
	client, err := capi.NewClient(cfg)
	check(err)
	check(client.Agent().ServiceDeregister(serviceID))
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func getArgs() (consulAddr string, err error) {
	consulAddr = os.Getenv("CONSUL_ADDR")
	if consulAddr == "" {
		err = errors.New("consul address is not set")
	}
	return
}

func healthcheck(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func main() {
	consulAddr, err := getArgs()
	check(err)
	serviceID, err := registerInConsul(consulAddr)
	check(err)
	defer unregisterConsul(consulAddr, serviceID)

	manager := AuthManager{}
	check(manager.loginManager.Initialize("auth-service", "pass", os.Getenv("POSTGRES_ADDR"), "5432"))
	check(manager.sessionsManager.Initialize(os.Getenv("HZ_CLUSTERNAME"), os.Getenv("HZ_MAP")))
	defer manager.Close()

	router := gin.Default()

	router.GET("/session_id", manager.InitializeSession)
	router.PATCH("/session_id", manager.LogIn)
	router.POST("/sign_up", manager.SingUp)
	router.GET("/coffee_shop_id", manager.ShopID)
	router.GET("/user_id", manager.UserID)
	router.GET("/healthcheck", healthcheck)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
