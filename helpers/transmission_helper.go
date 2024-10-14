package helpers

import (
	"encoding/json"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/nats-io/nats.go"
	"log"
	"os"
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

func ConnectToTransmission(msg *nats.Msg) (*transmissionrpc.Client, CommandNatsMessage) {
	log.Printf("[ConnectToTransmission] Отримано повідомлення з NATS: %s", string(msg.Data))

	var inputMessage CommandNatsMessage
	if err := json.Unmarshal(msg.Data, &inputMessage); err != nil {
		log.Printf("[ConnectToTransmission] Помилка при розборі повідомлення з NATS: %v", err)
		return nil, inputMessage
	}

	log.Printf("[ConnectToTransmission] Парсинг повідомлення: chatID = %d, arguments = %s", inputMessage.ChatId, inputMessage.Arguments)

	if inputMessage.ChatId != 0 {
		ip := os.Getenv("TORRENT_IP")
		port, err := strconv.ParseInt(os.Getenv("TORRENT_PORT"), 10, 64)
		if err != nil {
			log.Printf("[ConnectToTransmission] Помилка, порт Transmission задано невірно: %v", err)
			return nil, inputMessage
		}

		client, err := transmissionrpc.New(ip, "", "", &transmissionrpc.AdvancedConfig{
			Port:  uint16(port),
			HTTPS: false,
		})
		if err != nil {
			log.Printf("[ConnectToTransmission] Помилка підключенні до transmission: %v", err)
			return nil, inputMessage
		}
		return client, inputMessage
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
