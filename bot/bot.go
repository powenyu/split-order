package bot

import (
	"fmt" //to print errors
	"strings"

	"github.com/bwmarrin/discordgo"         //discordgo package from the repo of bwmarrin .
	"github.com/powenyu/split-order/config" //importing our config package which we have created above
)

var BotId string
var goBot *discordgo.Session

func Start() {

	//creating new bot session
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Making our bot a user using User function .
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Storing our id from u to BotId .
	BotId = u.ID

	// Adding handler function to handle our messages using AddHandler from discordgo package. We will declare messageHandler function later.
	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running !!")
}

//Definition of messageHandler function it takes two arguments first one is discordgo.Session which is s , second one is discordgo.MessageCreate which is m.
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check.
	if m.Author.ID == BotId {
		return
	}

	//If we message ping to our bot in our discord it will return us pong .
	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "pong")
		if err != nil {
			fmt.Println("send message, err: ", err.Error())
			return
		}
	}

	if len(m.Content) < 2 || m.Content[0:1] != config.BotPrefix {
		return
	}

	cmdline := strings.Split(m.Content, " ")

	switch cmdline[0] {
	case "!create":
		order, err := CreateOrder(m, s)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, err.Error()+"  你是不是不會打字？？")
			return
		}

		if !order.IsValid() {
			_, err := s.ChannelMessageSend(m.ChannelID, "發生問題ˋˊ 請私訊<@428929512193916928>釐清責任歸屬")
			if err != nil {
				fmt.Println("send message, err: ", err.Error())
			}
			return
		}

		//TODO: save order to database
		if err := order.Create(); err != nil {
			_, sendErr := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s <@428929512193916928>決定要去搖飲料了", err.Error()))
			if err != nil {
				fmt.Println("send message, err: ", sendErr.Error())
			}
			return
		}

		_, err = s.ChannelMessageSend(m.ChannelID, "資料插入成功 大概吧")
		if err != nil {
			fmt.Println("send message, err: ", err.Error())
			return
		}
		return
	case "!list":
		msg, err := List(m, s)
		if err != nil {
			_, sendErr := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s <@428929512193916928>決定要去搖飲料了", err.Error()))
			if err != nil {
				fmt.Println("send message, err: ", sendErr.Error())
			}
			return
		}

		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		if err != nil {
			fmt.Println("send message, err: ", err.Error())
		}
		return
	default:
		return
	}
}
