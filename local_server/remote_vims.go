package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"dev.local/types"
)

func get_remote_vms(local_port uint) (*[]types.ItermProfile, error) {
	var r []types.ItermProfile
	url := fmt.Sprintf(`http://localhost:%d/%s`, VIM_LOCAL_PORT, `list`)

	resp, err := http.Get(url)
	F(err)
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	F(err)

	fmt.Println(
		fmt.Sprintf("%s", bytes),
	)
  F(json.Unmarshal(bytes, &r))

	return &r, nil

}
