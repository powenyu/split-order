package config

import (
	"encoding/json"
	"fmt"       //used to print errors majorly.
	"io/ioutil" //it will be used to help us read our config.json file.
	"os"
)

var (
	Token        string
	BotPrefix    string
	Port         string
	DatabaseURL  string
	ReadTimeout  int = 180
	WriteTimeout int = 60
)

type configStruct struct {
	Token       string `json:"Token"`
	BotPrefix   string `json:"BotPrefix"`
	Port        string `json:"PORT"`
	DatabaseURL string `json:"DATABASE_URL"`
}

func ReadConfig() error {

	env := os.Getenv("env")
	fmt.Println("debug log : ", env)
	if env == "production" {
		Token = os.Getenv("Token")
		BotPrefix = os.Getenv("BotPrefix")
		Port = os.Getenv("PORT")
		DatabaseURL = os.Getenv("DATABASE_URL")
		return nil
	}

	file, err := ioutil.ReadFile("./env/config.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var config configStruct
	fmt.Println(string(file))
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix
	Port = config.Port
	DatabaseURL = config.DatabaseURL

	return nil
}
