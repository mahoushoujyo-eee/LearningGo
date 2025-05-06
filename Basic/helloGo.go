package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type requestBody struct {
	Key    string `json:"key"`
	Info   string `json:"info"`
	UserId string `json:"userId"`
}

type responseBody struct {
	Code int      `json:"code"`
	Text string   `json:"text"`
	List []string `json:"list"`
	Url  string   `json:"url"`
}

func process(inputChan <-chan string, userid string) {
	for {
		input := <-inputChan

		if input == "EOF" {
			break
		}

		reqData := &requestBody{
			Key:    "792bcf45156d488c92e9d11da494b085",
			Info:   input,
			UserId: userid,
		}

		byteData, _ := json.Marshal(&reqData)

		req, err := http.NewRequest("POST", "https://www.tuling123.com/openai/api", bytes.NewBuffer(byteData))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			var respData responseBody
			json.Unmarshal(body, &respData)
			fmt.Println("AI:", respData.Text)
		}

		if resp != nil {
			resp.Body.Close()
		}
	}
}

func main() {
	var input string
	fmt.Println("请输入内容，输入EOF结束：")
	channel := make(chan string)

	defer close(channel)

	go process(channel, "123456")
	for {
		fmt.Scanln(&input)
		channel <- input
		if input == "EOF" {
			channel <- input
			break
		}
	}
	fmt.Println("程序结束")
}
