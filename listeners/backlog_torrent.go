package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"strconv"
	"strings"
)

func ExecuteBacklogTorrent() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[ExecuteBacklogTorrent] Старт перенесення торенту у сховище для подивитися пізніше")

		return MoveTorrent(args, helpers.BacklogDir, client)
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_BACKLOG", processor)
}

func ExecuteDeBacklogTorrent() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[ExecuteDeBacklogTorrent] Старт перенесення торенту зі сховище для подивитися пізніше")

		return MoveTorrent(args, helpers.DownloadDir, client)
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_DEBACKLOG", processor)
}

func MoveTorrent(args []string, newLocation string, client *transmissionrpc.Client) string {
	strId := args[0]
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		text := fmt.Sprintf("[MoveTorrent] ID торента \"%s\" не валідний: %v", strId, err)
		log.Printf(text)
		return text
	}

	ctx := context.Background()
	err = client.TorrentSetLocation(ctx, id, newLocation, true)

	answer := fmt.Sprintf("[MoveTorrent] Почато операцію переміщення торанта ID=\"%d\"", id)

	switch {
	case err != nil && (ctx.Err() == context.DeadlineExceeded ||
		strings.Contains(err.Error(), "context deadline exceeded") ||
		strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers")):
	case err != nil:
		answer = fmt.Sprintf("[MoveTorrent] Помилка при зміні локації медіафайлів торента ID=\"%d\": %v", id, err)
		log.Printf(answer)
	default:
	}

	log.Printf(answer)
	return answer
}
