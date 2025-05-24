package helpers

import (
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"github.com/hekmon/transmissionrpc/v2"
	"log"
	"os"
	"strconv"
)

const DownloadDir = "/downloads/complete"
const BacklogDir = "/downloads/backlog"

func ListenToNatsMessages(queue string, f func(args []string, client *transmissionrpc.Client) []string) {
	processor := func(data []byte) {
		userId, args, err := nats_helper.ParseNatsBotCommand(data)
		if err != nil {
			log.Printf("[ListenToNatsMessages] Помилка під час прослуховування черги \"%s\" NATS: %v", queue, err)
			return
		}
		executeUserCommand(userId, args, f)
	}
	listener := &nats_helper.NatsListenerHandler{
		Function: processor,
	}
	err := nats_helper.StartNatsListener(queue, listener)
	if err != nil {
		log.Printf("[ListenToNatsMessages] Помилка під час прослуховування черги \"%s\" NATS: %v", queue, err)
	}
}

func executeUserCommand(userId int64, args []string, f func(args []string, client *transmissionrpc.Client) []string) {
	if userId == 0 {
		log.Printf("[ConnectToTransmission] Помилка: ID користувача порожній")
	}

	client := connectToTransmission()
	result := f(args, client)

	answer := ""
	sendLast := false
	for i := 0; i < len(result); i++ {
		if len(answer)+len(result[i]) > 4000 {
			nats_helper.SendMessageToUser(userId, answer)
			answer = ""
			sendLast = false
		} else {
			answer = answer + "\n" + result[i]
			sendLast = true
		}
	}
	if sendLast {
		nats_helper.SendMessageToUser(userId, answer)
	}
}

func connectToTransmission() *transmissionrpc.Client {
	portStr := os.Getenv("TORRENT_PORT")
	torrentIp := os.Getenv("TORRENT_IP")
	port, err := strconv.Atoi(portStr) // Преобразуем строку в int
	if err != nil {
		log.Fatalf("[ConnectToTransmission] Invalid port value in TORRENT_PORT: %v", err)
		return nil
	}

	if port < 0 || port > 65535 {
		log.Fatalf("[ConnectToTransmission] Invalid port range in TORRENT_PORT: %d. Must be between 0 and 65535.", port)
		return nil
	}

	client, err := transmissionrpc.New(torrentIp, "", "", &transmissionrpc.AdvancedConfig{Port: uint16(port), HTTPS: false})
	if err != nil {
		log.Printf("[ConnectToTransmission] Помилка при підключенні до transmission: %v", err)
		return nil
	}

	return client
}
