package postgresql

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/powenyu/split-order/config"
)

var (
	Pool *pgxpool.Pool
)

func Initialize() {
	databaseURL := config.DatabaseURL
	u, err := url.Parse(databaseURL)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Add("target_session_attrs", "read-write")
	q.Add("connect_timeout", "10")
	q.Add("pool_max_conns", "50")
	q.Add("pool_max_conn_lifetime", "180s")
	q.Add("pool_max_conn_idle_time", "180s")
	u.RawQuery = q.Encode()
	databaseURL = u.String()

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		panic(err)
	}

	Pool = setupPool(cfg)
}

func setupPool(cfg *pgxpool.Config) *pgxpool.Pool {
	log.Printf("[info] Connecting to Postgresql @ %v", cfg.ConnConfig.ConnString())
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}
	log.Printf("[info] Postgresql initialization is done")
	return pool
}

func Dispose() {
	Pool.Close()
}
