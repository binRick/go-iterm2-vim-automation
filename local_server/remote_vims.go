package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func get_remote_vms(local_port uint) (*[]ItermProfile, error) {
	var r []ItermProfile
	url := fmt.Sprintf(`http://localhost:%d/%s`, VIM_LOCAL_PORT, `list`)

	resp, err := http.Get(url)
	F(err)
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	F(err)

	fmt.Println(
		fmt.Sprintf("%s", bytes),
	)
	err = json.Unmarshal(bytes, &r)
	F(err)

	return &r, nil

}
