package config

import (
	"encoding/json"
	"fmt"       //used to print errors majorly.
	"io/ioutil" //it will be used to help us read our config.json file.
	"os"
)

var (
	Token     string //To store value of Token from config.json .
	BotPrefix string // To store value of BotPrefix from config.json.

	config *configStruct //To store value extracted from config.json.
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {

	env := os.Getenv("env")
	fmt.Println("debug log : ", env)
	if env == "production" {
		Token = os.Getenv("Token")
		BotPrefix = os.Getenv("BotPrefix")

		return nil
	}

	file, err := ioutil.ReadFile("./env/config.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}
