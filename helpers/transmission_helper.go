package helpers

import (
	"encoding/json"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"github.com/halushko/kino-cat-core-go/warehouse_helper"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/nats-io/nats.go"
	"log"
	"sort"
	"strconv"
)

type TelegramUserNatsMessage struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type CommandNatsMessage struct {
	ChatId    int64    `json:"chat_id"`
	Arguments []string `json:"arguments"`
}

const OutputQueue = "TELEGRAM_OUTPUT_TEXT_QUEUE"

func executeForServers(msg *nats.Msg, f func(key string, args []string, client *transmissionrpc.Client) string) {
	clients, inputMessage := connectToTransmission(msg)
	keys := make([]string, 0)
	for key := range clients {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, key := range keys {
		log.Printf("[ExecuteForServers] Старт роботи зі сховищем %s", key)
		SendAnswer(inputMessage.ChatId, f(key, inputMessage.Arguments, clients[key]))
	}
}

func ListenToNatsMessages(queue string, f func(key string, args []string, client *transmissionrpc.Client) string) {
	processor := func(msg *nats.Msg) {
		executeForServers(msg, f)
	}
	listener := &nats_helper.NatsListenerHandler{
		Function: processor,
	}
	err := nats_helper.StartNatsListener(queue, listener)
	if err != nil {
		log.Printf("[ListenToNatsMessages] Помилка під час прослуховування черги \"%s\" NATS: %v", queue, err)
	}
}

func connectToTransmission(msg *nats.Msg) (map[string]*transmissionrpc.Client, CommandNatsMessage) {
	log.Printf("[ConnectToTransmission] Отримано повідомлення з NATS: %s", string(msg.Data))

	var inputMessage CommandNatsMessage
	if err := json.Unmarshal(msg.Data, &inputMessage); err != nil {
		log.Printf("[ConnectToTransmission] Помилка при розборі повідомлення з NATS: %v", err)
		return nil, inputMessage
	}

	log.Printf("[ConnectToTransmission] Парсинг повідомлення: chatID = %d, arguments = %s", inputMessage.ChatId, inputMessage.Arguments)

	if inputMessage.ChatId != 0 {
		var clients = make(map[string]*transmissionrpc.Client)
		for key, value := range getTransmissionServers(inputMessage.Arguments) {
			ip := value.IP
			port, err := strconv.ParseInt(value.Port, 10, 64)

			if err != nil {
				log.Printf("[ConnectToTransmission] Помилка, порт Transmission задано невірно для сховища %s %v", key, err)
				continue
			}

			client, err := transmissionrpc.New(ip, "", "", &transmissionrpc.AdvancedConfig{
				Port:  uint16(port),
				HTTPS: false,
			})
			if err != nil {
				log.Printf("[ConnectToTransmission] Помилка підключенні до transmission для сховища %s: %v", key, err)
				continue
			}
			clients[key] = client
		}
		return clients, inputMessage
	} else {
		log.Printf("[ConnectToTransmission] Помилка: ID користувача порожній")
	}
	return nil, inputMessage
}

func SendAnswer(chatId int64, message string) {
	log.Printf("[SendAnswer] Answer:%s", message)
	if request, errMarshal := json.Marshal(TelegramUserNatsMessage{
		ChatId: chatId,
		Text:   message,
	}); errMarshal == nil {
		if errPublish := nats_helper.PublishToNATS(OutputQueue, request); errPublish != nil {
			log.Printf("[SendAnswer] ERROR in publish to %s:%s", OutputQueue, errPublish)
		}
	} else {
		log.Printf("[SendAnswer] ERROR in publish to %s:%s", OutputQueue, errMarshal)
	}
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
