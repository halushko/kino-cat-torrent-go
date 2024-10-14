package listeners

import (
	"context"
	"fmt"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"github.com/nats-io/nats.go"
	"kino-cat-torrent-go/helpers"
	"log"
	"strconv"
)

func ExecuteResumeTorrent() {
	processor := func(msg *nats.Msg) {
		client, inputMessage := helpers.ConnectToTransmission(msg)

		log.Printf("[ExecuteResumeTorrent] Старт зупинки торенту")
		strId := inputMessage.Arguments[len(inputMessage.Arguments)-1]
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			log.Printf("[ExecuteResumeTorrent] ID торента \"%s\" не валідний: %v", strId, err)
			return
		}

		err = client.TorrentStopIDs(context.Background(), []int64{id})

		answer := ""
		if err != nil {
			answer = fmt.Sprintf("Торент з ID=%d не поновлено", id)
		} else {
			answer = fmt.Sprintf("Торент з ID=%d поновлено", id)
		}

		helpers.SendAnswer(inputMessage.ChatId, answer)
	}

	listener := &nats_helper.NatsListener{
		Handler: processor,
	}

	nats_helper.StartNatsListener("EXECUTE_TORRENT_COMMAND_RESUME_TORRENT", listener)
}
