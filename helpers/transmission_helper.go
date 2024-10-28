package helpers

import (
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"github.com/halushko/kino-cat-core-go/warehouse_helper"
	"github.com/hekmon/transmissionrpc/v2"
	"log"
	"sort"
	"strconv"
)

const outputQueue = "TELEGRAM_OUTPUT_TEXT_QUEUE"

func ListenToNatsMessages(queue string, f func(key string, args []string, client *transmissionrpc.Client) string) {
	processor := func(data []byte) {
		userId, args, err := nats_helper.ParseNatsBotCommand(data)
		if err != nil {
			log.Printf("[ListenToNatsMessages] Помилка під час прослуховування черги \"%s\" NATS: %v", queue, err)
			return
		}
		executeForServers(userId, args, f)
	}
	listener := &nats_helper.NatsListenerHandler{
		Function: processor,
	}
	err := nats_helper.StartNatsListener(queue, listener)
	if err != nil {
		log.Printf("[ListenToNatsMessages] Помилка під час прослуховування черги \"%s\" NATS: %v", queue, err)
	}
}

func SendAnswer(userId int64, message string) {
	log.Printf("[SendAnswer] Відповідь: %s", message)
	err := nats_helper.PublishTextMessage(outputQueue, userId, message)
	if err != nil {
		log.Printf("[SendAnswer] Помилка: %v", err)
	}
}

func executeForServers(userId int64, args []string, f func(key string, args []string, client *transmissionrpc.Client) string) {
	clients := connectToTransmission(userId, args)

	warehouses := make([]string, 0)
	for key := range clients {
		warehouses = append(warehouses, key)
	}
	sort.Slice(warehouses, func(i, j int) bool { return warehouses[i] < warehouses[j] })

	for _, key := range warehouses {
		log.Printf("[ExecuteForServers] Старт роботи зі сховищем %s", key)
		SendAnswer(userId, f(key, args, clients[key]))
	}
}

func connectToTransmission(userId int64, args []string) map[string]*transmissionrpc.Client {
	if userId == 0 {
		log.Printf("[ConnectToTransmission] Помилка: ID користувача порожній")
		return nil
	}

	var clients = make(map[string]*transmissionrpc.Client)
	for key, value := range getTransmissionServers(args) {
		port, err := strconv.ParseInt(value.Port, 10, 64)
		if err != nil {
			log.Printf("[ConnectToTransmission] Помилка, порт Transmission задано невірно для сховища %s %v", key, err)
			continue
		}

		client, err := transmissionrpc.New(value.IP, "", "", &transmissionrpc.AdvancedConfig{Port: uint16(port), HTTPS: false})
		if err != nil {
			log.Printf("[ConnectToTransmission] Помилка підключенні до transmission для сховища %s: %v", key, err)
			continue
		}

		clients[key] = client
	}
	return clients
}

func getTransmissionServers(commandArguments []string) map[string]warehouse_helper.WarehouseConfig {
	warehouses, err := warehouse_helper.ParseWarehouseConfig()
	if err != nil {
		log.Printf("[getTransmissionServers] Can't cet transmission servers: %v", err)
		return nil
	}
	result := make(map[string]warehouse_helper.WarehouseConfig)

	if len(warehouses) == 1 {
		for _, config := range warehouses {
			result[config.Name] = config
		}
	} else if len(commandArguments) == 0 {
		for _, config := range warehouses {
			result[config.Name] = config
		}
	} else {
		_, exists := warehouses[commandArguments[0]]
		if !exists {
			for _, config := range warehouses {
				result[config.Name] = config
			}
		} else {
			for _, config := range warehouses {
				result[config.Name] = config
			}
		}
	}
	return result
}
