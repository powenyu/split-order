package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	postgresql "github.com/powenyu/split-order/postgres"
)

type Order struct {
	ID          uint      `json:"id" db:"id"`
	Collective  string    `json:"collective" db:"collective"`
	GuildID     string    `json:"guild_id" db:"guild_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Description string    `json:"description" db:"description"`

	//extend
	OrderParticipants []OrderParticipant `json:"order_participants"`
}

func (o *Order) Create() error {
	values := make([]string, 0, len(o.OrderParticipants))
	for _, v := range o.OrderParticipants {
		value := fmt.Sprintf("((SELECT id FROM t1),'%s',%f)", v.UserID, v.Price)
		values = append(values, value)
	}

	sql := fmt.Sprintf(`
	WITH t1 AS(
		INSERT INTO
			orders
		(collective, guild_id, created_at, description)
		VALUES
		($1,$2,$3,$4)
		RETURNING id
	)
	
	INSERT INTO
		order_participants
	(order_id,user_id,price)
	VALUES
	%s
	`, strings.Join(values, ","))

	if _, err := postgresql.Pool.Exec(context.Background(), sql, o.Collective, o.GuildID, o.CreatedAt, o.Description); err != nil {
		return err
	}
	return nil
}

func (o *Order) IsValid() bool {
	//check total
	if len(o.OrderParticipants) < 1 || o.Collective == "" {
		return false
	}

	var subtotal float64
	var pay float64
	for _, v := range o.OrderParticipants {
		if v.Price < 0 {
			pay += v.Price
		} else {
			subtotal += v.Price
		}
	}
	if subtotal >= pay*-1 {
		return false
	}

	return true
}

func SelectDistinctCollective() ([]string, error) {
	var collectives []string

	sql := `
	SELECT
		distinct(collective)
	FROM
		orders
	`

	rows, err := postgresql.Pool.Query(context.Background(), sql)
	if err != nil {
		return collectives, err
	}

	if err := pgxscan.ScanAll(&collectives, rows); err != nil {
		return collectives, err
	}

	return collectives, nil
}

func SelectCollectiveRecord(collective string) (*[]Order, error) {
	var orders []Order

	sql := `
	SELECT
		collective,
		guild_id,
		created_at,
		description,
		JSON_AGG(JSON_BUILD_OBJECT(
			'user_id',op.user_id,
			'price', op.price::numeric
		)) AS order_participants
	FROM
		orders AS o
	LEFT JOIN(
		SELECT
			*
		FROM
			order_participants
	) AS op ON o.id = op.order_id
	WHERE
		o.collective = $1
	GROUP BY
		collective,guild_id,created_at,description
	`

	rows, err := postgresql.Pool.Query(context.Background(), sql, collective)
	if err != nil {
		return &orders, err
	}

	if err := pgxscan.ScanAll(&orders, rows); err != nil {
		return &orders, err
	}

	return &orders, nil
}
