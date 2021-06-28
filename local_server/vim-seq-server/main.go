package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"mrz.io/itermctl"

	"github.com/k0kubun/pp"
)

func main() {
	fmt.Println("vim-go")
}

func monitor_control_seq() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	re := regexp.MustCompile("^test-seq:.*")
	fmt.Println(
		"CONTROL_SEQUENCE_NAME:", CONTROL_SEQUENCE_NAME,
	)
	notifications, err := itermctl.MonitorCustomControlSequences(ctx, _conn, CONTROL_SEQUENCE_NAME, re, itermctl.AllSessions)
	F(err)

	msg := fmt.Sprintf(`
 ** Waiting for control sequeunces **

CONTROL_SEQUENCE_NAME: %v

`,
		CONTROL_SEQUENCE_NAME,
	)
	pp.Println(msg)
	dm := func() {
		select {
		case notification := <-notifications:
			pp.Println("New Sequence!")
			pp.Println(fmt.Sprintf(`Session: %s`, notification.Notification.GetSession()))
			pp.Println(fmt.Sprintf(`Matches qty: %d`, len(notification.Matches)))
			for _, m := range notification.Matches {
				pp.Println(fmt.Sprintf(`   Match:     %d bytes`, len(m)))
				if len(strings.Split(m, `:`)) == 2 {
					seq_json := map[string]interface{}{}
					seq_enc := strings.Split(m, `:`)[1]

					seq_dec, err := base64.StdEncoding.DecodeString(seq_enc)
					F(err)

					seq_dec_trimmed := strings.Trim(fmt.Sprintf(`%s`, seq_dec), "")

					dec_err := json.Unmarshal([]byte(seq_dec_trimmed), &seq_json)
					F(dec_err)

					pp.Println(
						seq_json,
					)
				}

			}
		}
	}

	for {
		dm()
	}
}

func F(err error) {
	if err != nil {
		log.Error(err)
		panic(err)
	}
}
