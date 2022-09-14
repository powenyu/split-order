package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/powenyu/split-order/postgres/model"
)

var (
	errInvalidCmd error = errors.New("Invalid params")
)

type userInfo struct {
	userID   string
	userName string
	price    float64
}

type rowInfo struct {
	userPrice   map[string]float64
	description string
	createdAt   time.Time
}

func PrettyPrint(msg string, obj interface{}) {
	log.Println(msg)
	s, _ := json.MarshalIndent(obj, "", "\t")
	log.Println(string(s))
}

func CreateOrder(m *discordgo.MessageCreate, s *discordgo.Session) (*model.Order, error) {
	var order model.Order
	var orderParticipants []model.OrderParticipant

	msg := m.Content

	cmdlines := strings.Split(msg, " ")
	for _, cmdline := range cmdlines {
		if strings.Contains(cmdline, ":") {
			orderParticipant, err := parseOrderParticipant(cmdline, s)
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

func parseOrderParticipant(msg string, s *discordgo.Session) (model.OrderParticipant, error) {
	participantString := strings.Split(msg, ":")
	var orderParticipant model.OrderParticipant

	//only allow one ":" so there will only be split to two parts
	if len(participantString) != 2 {
		fmt.Println("participate error")
		return orderParticipant, errInvalidCmd
	}

	_, err := getUser(participantString[0], s)
	if err != nil {
		return orderParticipant, err
	}

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

func List(m *discordgo.MessageCreate, s *discordgo.Session) (string, error) {
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

		comment, err = drawDiagram(*orders, s)
		if err != nil {
			return comment, err
		}

	}
	return comment, nil
}

//+ , - . 0 ♦ # ° ± n ↓ ┘ ┐ ┌ └ ┼ ⎺ ⎻ ─ ⎼ ⎽ ├ ┤ ┴ ┬ ≤ │ ≥ # ≠ £ ·
func drawDiagram(orders []model.Order, s *discordgo.Session) (string, error) {
	var comment string

	PrettyPrint("orders", orders)

	allUsers := make(map[string]string)
	allPrice := make(map[string]float64)
	userPrice := make(map[string]float64)
	userRows := make([]rowInfo, 0, len(orders))
	for _, order := range orders {
		for _, participant := range order.OrderParticipants {
			_, ok := allUsers[participant.UserID]
			if !ok {
				user, err := getUser(participant.UserID, s)
				if err != nil {
					return comment, err
				}
				allUsers[participant.UserID] = user.Username
			}

			userPrice[participant.UserID] += participant.Price
		}

		count(userPrice)
		PrettyPrint("userPrice2", userPrice)
		for k, v := range userPrice {
			allPrice[k] += v
		}
		userRow := rowInfo{
			userPrice:   userPrice,
			createdAt:   order.CreatedAt,
			description: order.Description,
		}
		userRows = append(userRows, userRow)
		userPrice = make(map[string]float64)
	}

	//set first row
	i := 0
	fixuser := make([]string, 0, len(allUsers))
	for k, v := range allUsers {
		if i < len(allUsers)-1 {
			comment += fmt.Sprintf("%-30s|", v)
		} else {
			comment += fmt.Sprintf("%-30s\n", v)
		}
		fixuser = append(fixuser, k)
		i++
	}

	for _, userRow := range userRows {
		i = 0
		for _, k := range fixuser {
			if i < len(allUsers)-1 {
				fmt.Print(k, ",")
				comment += fmt.Sprintf("%-30s|", fmt.Sprint(int(userRow.userPrice[k])))
			} else {
				fmt.Print(k, "\n")
				comment += fmt.Sprintf("%-30s\n", fmt.Sprint(int(userRow.userPrice[k])))
			}
			i++
		}
	}

	comment += "\n"
	i = 0
	for _, k := range fixuser {
		if i < len(allUsers)-1 {
			fmt.Print(k, ",")
			comment += fmt.Sprintf("%-30d|", int(allPrice[k]))
		} else {
			fmt.Print(k, "\n")
			comment += fmt.Sprintf("%-30d\n", int(allPrice[k]))
		}
		i++
	}

	fmt.Print(comment)

	return comment, nil
}

func getUser(userID string, s *discordgo.Session) (*discordgo.User, error) {
	user, err := s.User(strings.Trim(userID, "<>@"))
	return user, err
}

// TODO: draw diagram to show detail
func draw(alluser map[string]userInfo, userRows []rowInfo) string {

	return ""
}

func count(rowprice map[string]float64) {
	var pay float64 = 0
	// count should pay
	for _, v := range rowprice {
		pay -= v
	}

	avg := pay / float64(len(rowprice))

	for k := range rowprice {
		rowprice[k] += avg
	}
}

/*
 description | user1 | user2 | user3 | create at
 ------------┼-------┼-------┼-------┼-----------
             | 300   | 400   | -700  |

*/
