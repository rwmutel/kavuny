package main

import (
	"errors"
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"net"
	"os"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
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
	return "", errors.New("are you connected to the network")
}

func registerInConsul(consulAddr string) (string, string, string, error) {
	cfg := capi.DefaultConfig()
	cfg.Address = consulAddr
	client, err := capi.NewClient(cfg)
	if err != nil {
		return "", "", "", err
	}
	host, err := externalIP()
	if err != nil {
		return "", "", "", err
	}
	serviceID, err := uuid.NewUUID()
	if err != nil {
		return "", "", "", err
	}
	err = client.Agent().ServiceRegister(&capi.AgentServiceRegistration{
		ID:      serviceID.String(),
		Name:    "log-service",
		Port:    8080,
		Address: host,
	})
	if err != nil {
		return "", "", "", err
	}
	keyValuePairs := client.KV()
	kafkaAddrKV, _, err := keyValuePairs.Get("kafka_address", &capi.QueryOptions{})
	if err != nil {
		return "", "", "", err
	}
	kafkaTopicKV, _, err := keyValuePairs.Get("kafka_topic", &capi.QueryOptions{})
	if err != nil {
		return "", "", "", err
	}
	kafkaAddr := string(kafkaAddrKV.Value)
	kafkaTopic := string(kafkaTopicKV.Value)
	return serviceID.String(), kafkaAddr, kafkaTopic, nil
}

func unregisterConsul(consulAddr string, serviceID string) {
	cfg := capi.DefaultConfig()
	cfg.Address = consulAddr
	client, err := capi.NewClient(cfg)
	check(err)
	check(client.Agent().ServiceDeregister(serviceID))
}

func getArgs() (consulAddr string, err error) {
	consulAddr = os.Getenv("CONSUL_ADDR")
	if consulAddr == "" {
		err = errors.New("consul address is not set")
	}
	return
}

func consume[T any](data chan<- string, stop <-chan T, consumer *kafka.Consumer) error {
	defer close(data)
	for {
		select {
		case <-stop:
			return nil
		default:
			message, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}
			str := string(message.Value)
			data <- str
		}
	}
}

func consumeFromKafka[T any](data chan<- string, stop <-chan T, address, topic string) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": address,
		"group.id":          "message",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return err
	}
	defer consumer.Close()
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		return err
	}
	return consume(data, stop, consumer)
}

func main() {
	consulAddr, err := getArgs()
	check(err)
	serviceID, kafkaAddr, kafkaTopic, err := registerInConsul(consulAddr)
	check(err)
	defer unregisterConsul(consulAddr, serviceID)

	data := make(chan string)
	stop := make(chan bool)

	go func() {
		check(consumeFromKafka(data, stop, kafkaAddr, kafkaTopic))
	}()

	logFile, err := os.Open(os.Getenv("LOG_FILE"))
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(logFile)
	defer logFile.Close()

	for logLine := range data {
		log.Println(logLine)
	}
}
