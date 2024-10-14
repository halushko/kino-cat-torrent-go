package listeners

import (
	"context"
	"encoding/json"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"strconv"
)

type CommandNatsMessage struct {
	ChatId    int64    `json:"chat_id"`
	Arguments []string `json:"arguments"`
}

func GetAllTorrents() {
	processor := func(msg *nats.Msg) {
		log.Printf("[GetAllTorrents] Отримано повідомлення з NATS: %s", string(msg.Data))

		var inputMessage CommandNatsMessage
		if err := json.Unmarshal(msg.Data, &inputMessage); err != nil {
			log.Printf("[GetAllTorrents] Помилка при розборі повідомлення з NATS: %v", err)
			return
		}

		log.Printf("[GetAllTorrents] Парсинг повідомлення: chatID = %d, arguments = %s", inputMessage.ChatId, inputMessage.Arguments)

		if inputMessage.ChatId != 0 {
			ip := os.Getenv("TORRENT_IP")
			port, err := strconv.ParseInt(os.Getenv("TORRENT_PORT"), 10, 64)
			if err != nil {
				log.Printf("[GetAllTorrents] Помилка, порт Transmission задано невірно: %v", err)
				return
			}

			client, err := transmissionrpc.New(ip, "", "", &transmissionrpc.AdvancedConfig{
				Port:  uint16(port),
				HTTPS: false,
			})
			if err != nil {
				log.Printf("[GetAllTorrents] Помилка підключенні до transmission: %v", err)
				return
			}

			torrents, err := client.TorrentGetAll(context.Background())
			if err != nil {
				log.Printf("[GetAllTorrents] Помилка отримання переліку торентов: %v", err)
				return
			}

			for _, torrent := range torrents {
				log.Printf("ID: %d, Назві: %s\n", *torrent.ID, *torrent.Name)
			}
		} else {
			log.Printf("[GetAllTorrents] Помилка: ID користувача чи текст повідомлення порожні")
		}
	}

	listener := &nats_helper.NatsListener{
		Handler: processor,
	}

	nats_helper.StartNatsListener("EXECUTE_TORRENT_COMMAND_LIST", listener)
}
