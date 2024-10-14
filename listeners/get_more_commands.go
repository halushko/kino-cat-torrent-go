package listeners

import (
	"context"
	"fmt"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/nats-io/nats.go"
	"kino-cat-torrent-go/helpers"
	"log"
	"strconv"
)

func GetMoreCommands() {
	processor := func(msg *nats.Msg) {
		client, inputMessage := helpers.ConnectToTransmission(msg)

		log.Printf("[GetMoreCommands] Старт отримання інформації по торенту")
		strId := inputMessage.Arguments[len(inputMessage.Arguments)-1]
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			log.Printf("[GetMoreCommands] ID торента \"%s\" не валідний: %v", strId, err)
			return
		}

		torrents, err := client.TorrentGet(context.Background(), []string{"name", "id", "status"}, []int64{id})
		if err != nil {
			log.Printf("[GetMoreCommands] Помилка отримання переліку торентов: %v", err)
			return
		}
		answer := ""
		if len(torrents) == 1 {
			log.Printf("[GetMoreCommands] Інформацію про торент \"%d\" отримано", id)
			answer = generateAnswerMore(torrents[0])
		} else {
			log.Printf("[GetMoreCommands] Інформації про торент \"%d\" немає", id)
			answer = fmt.Sprintf("Нажаль для торента з ID=%d не можна отримати Ім'я", id)
		}
		helpers.SendAnswer(inputMessage.ChatId, answer)
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
