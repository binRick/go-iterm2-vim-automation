package main

import (
	"github.com/manifoldco/promptui"

	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"mrz.io/itermctl"
	"github.com/c-bata/go-prompt"

	"github.com/k0kubun/pp"
)

const (
SEQ_PREFIX = `test-seq`
)

func main() {
	conn, err := itermctl.GetCredentialsAndConnect("itermctl_statusbar_example", true)
	F(err)
	_conn = conn
	app, err := itermctl.NewApp(_conn)
	F(err)
	_app = app
	go monitor_control_seq()
	for {
		time.Sleep(5 * time.Second)
	}
}
func pui() {
items := []string{"Vim", "Emacs", "Sublime", "VSCode", "Atom"}
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    "What's your text editor",
			Items:    items,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %s\n", result)
}
func monitor_control_seq() {
pui()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	re := regexp.MustCompile(fmt.Sprintf("^%s:.*",SEQ_PREFIX))
	notifications, err := itermctl.MonitorCustomControlSequences(ctx, _conn, CONTROL_SEQUENCE_NAME, re, itermctl.AllSessions)
	F(err)

	msg := fmt.Sprintf(`
 ** Waiting for control sequeunces **

CONTROL_SEQUENCE_NAME: %v
SEQ_PREFIX: %v

`,
		CONTROL_SEQUENCE_NAME,
		SEQ_PREFIX,
	)
	fmt.Println(msg)
	dm := func() {
		select {
		case notification := <-notifications:
			for _, m := range notification.Matches {
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
