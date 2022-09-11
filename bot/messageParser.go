package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/powenyu/split-order/postgres/model"
)

var (
	errInvalidCmd error = errors.New("Invalid params")
)

func PrettyPrint(msg string, obj interface{}) {
	log.Println(msg)
	s, _ := json.MarshalIndent(obj, "", "\t")
	log.Println(string(s))
}

func CreateOrder(m *discordgo.MessageCreate) (*model.Order, error) {
	var order model.Order
	var orderParticipants []model.OrderParticipant

	msg := m.Content

	cmdlines := strings.Split(msg, " ")
	for _, cmdline := range cmdlines {
		if strings.Contains(cmdline, ":") {
			orderParticipant, err := parseOrderParticipant(cmdline)
			if err != nil {
				fmt.Println("parse error")
				return &order, err
			}

			orderParticipants = append(orderParticipants, orderParticipant)
		} else if strings.Contains(cmdline, "-d") {
			comment := strings.TrimPrefix(cmdline, "-d=")
			comment = strings.Trim(comment, "\"")

			if comment == "" {
				fmt.Println("comment error")
				return &order, errInvalidCmd
			}

			order.Description = comment
		} else if strings.Contains(cmdline, "-c") {
			collective := strings.TrimPrefix(cmdline, "-c=")

			if collective == "" {
				fmt.Println("collective error")
				return &order, errInvalidCmd
			}

			order.Collective = collective
		}
	}

	order.OrderParticipants = orderParticipants
	order.GuildID = m.GuildID
	order.CreatedAt = m.Timestamp

	return &order, nil
}

func parseOrderParticipant(msg string) (model.OrderParticipant, error) {
	participantString := strings.Split(msg, ":")
	var orderParticipant model.OrderParticipant

	//only allow one ":" so there will only be split to two parts
	if len(participantString) != 2 {
		fmt.Println("participate error")
		return orderParticipant, errInvalidCmd
	}

	//TODO : check valid userid

	//check valid price
	price, err := strconv.ParseFloat(participantString[1], 64)
	if err != nil {
		fmt.Println("invalid price")
		return orderParticipant, errInvalidCmd
	}

	orderParticipant.Price = price
	orderParticipant.UserID = participantString[0]

	return orderParticipant, nil
}

func parseComment(msg string) (string, error) {
	var comment string

	comments := strings.Split(msg, "\"")

	fmt.Println("comments", comments)

	//need comment between two " so there will split to three parts
	if len(comment) != 3 {
		return comment, errInvalidCmd
	}

	//get buttom between "
	comment = comments[1]

	return comment, nil
}

func List(m *discordgo.MessageCreate) (string, error) {
	var comment string
	msg := m.Content

	cmdlines := strings.Split(msg, " ")

	if msg == "!list" {
		collectives, err := model.SelectDistinctCollective()
		if err != nil {
			return comment, err
		}

		if len(collectives) == 0 {
			return "NULL", nil
		}
		for i, v := range collectives {
			comment += fmt.Sprintf("%d: %s \n", i+1, v)
		}
	} else if len(cmdlines) == 2 {
		collective := cmdlines[1]
		orders, err := model.SelectCollectiveRecord(collective)
		if err != nil {
			return comment, err
		}

		comment, err = drawDiagram(*orders)
		if err != nil {
			return comment, err
		}

	}
	return comment, nil
}

//+ , - . 0 ♦ # ° ± n ↓ ┘ ┐ ┌ └ ┼ ⎺ ⎻ ─ ⎼ ⎽ ├ ┤ ┴ ┬ ≤ │ ≥ # ≠ £ ·
func drawDiagram(o []model.Order) (string, error) {
	var comment string

	// row := len(o)
	// var rows []int

	// for _, order := range o {
	// 	for _, user := range order.OrderParticipants {

	// 	}
	// }

	return comment, nil
}
