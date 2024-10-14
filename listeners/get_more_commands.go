package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"strconv"
)

func GetMoreCommands() {
	processor := func(msg *nats.Msg) {
		log.Printf("[GetMoreCommands] Отримано повідомлення з NATS: %s", string(msg.Data))

		var inputMessage CommandNatsMessage
		if err := json.Unmarshal(msg.Data, &inputMessage); err != nil {
			log.Printf("[GetMoreCommands] Помилка при розборі повідомлення з NATS: %v", err)
			return
		}

		log.Printf("[GetMoreCommands] Парсинг повідомлення: chatID = %d, arguments = %s", inputMessage.ChatId, inputMessage.Arguments)

		if inputMessage.ChatId != 0 {
			ip := os.Getenv("TORRENT_IP")
			port, err := strconv.ParseInt(os.Getenv("TORRENT_PORT"), 10, 64)
			if err != nil {
				log.Printf("[GetMoreCommands] Помилка, порт Transmission задано невірно: %v", err)
				return
			}

			client, err := transmissionrpc.New(ip, "", "", &transmissionrpc.AdvancedConfig{
				Port:  uint16(port),
				HTTPS: false,
			})
			if err != nil {
				log.Printf("[GetMoreCommands] Помилка підключенні до transmission: %v", err)
				return
			}
			log.Printf("[GetMoreCommands] Старт отримання інформації по торенту")
			strId := inputMessage.Arguments[len(inputMessage.Arguments)-1]
			id, err := strconv.ParseInt(strId, 10, 64)
			if err != nil {
				log.Printf("[GetMoreCommands] ID торента \"%s\" не валідний: %v", strId, err)
				return
			}

			torrents, err := client.TorrentGet(context.Background(), []string{"name", "id"}, []int64{id})
			if err != nil {
				log.Printf("[GetMoreCommands] Помилка отримання переліку торентов: %v", err)
				return
			}
			queue := "TELEGRAM_OUTPUT_TEXT_QUEUE"
			answer := ""
			if len(torrents) == 1 {
				log.Printf("[GetMoreCommands] Інформацію про торент \"%d\" отримано", id)
				answer = generateAnswerMore(torrents[0])
			} else {
				log.Printf("[GetMoreCommands] Інформації про торент \"%d\" немає", id)
				answer = fmt.Sprintf("Нажаль для торента з ID=%d не можна отримати Ім'я", id)
			}
			log.Printf("[GetMoreCommands] Answer:%s", answer)
			if request, errMarshal := json.Marshal(TelegramUserNatsMessage{
				ChatId: inputMessage.ChatId,
				Text:   answer,
			}); errMarshal == nil {
				if errPublish := nats_helper.PublishToNATS(queue, request); errPublish != nil {
					log.Printf("[GetMoreCommands] ERROR in publish to %s:%s", queue, errPublish)
				}
			} else {
				log.Printf("[GetMoreCommands] ERROR in publish to %s:%s", queue, errMarshal)
			}
		} else {
			log.Printf("[GetMoreCommands] Помилка: ID користувача порожній")
		}
	}

	listener := &nats_helper.NatsListener{
		Handler: processor,
	}

	nats_helper.StartNatsListener("EXECUTE_TORRENT_COMMAND_SHOW_COMMANDS", listener)
}

func generateAnswerMore(torrent transmissionrpc.Torrent) string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		*torrent.Name,
		fmt.Sprintf("%s%s%s%d", "/pause_", "", "", *torrent.ID),
		fmt.Sprintf("%s%s%s%d", "/resume_", "", "", *torrent.ID),
		fmt.Sprintf("%s%s%s%d", "/info_", "", "", *torrent.ID),
		fmt.Sprintf("%s%s%s%d", "/remove_", "", "", *torrent.ID),
	)
}
