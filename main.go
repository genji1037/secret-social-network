package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	_ "net/http/pprof"
	"secret-social-network/app/config"
	"secret-social-network/app/dgraph"
	"secret-social-network/app/server"
	"secret-social-network/app/storage"
)

func main() {
	// load config
	err := config.LoadConfig("config/server.yml")
	if err != nil {
		panic(err)
	}
	cfg := config.GetServe()

	log.StandardLogger().SetLevel(cfg.Level.Value())

	// TODO:⭐ init module parallel
	// conn mysql
	err = storage.Open(storage.Connection{
		Host:         cfg.MySQL.Host,
		User:         cfg.MySQL.User,
		Password:     cfg.MySQL.Password,
		Database:     cfg.MySQL.Database,
		Charset:      cfg.MySQL.Charset,
		MaxIdleConns: cfg.MySQL.MaxIdleConns,
		MaxOpenConns: cfg.MySQL.MaxOpenConns,
	})
	if err != nil {
		log.Panicf("[MAIN] failed to open mysql: %s", err.Error())
	}

	// conn dGraph
	if err = dgraph.Open(cfg.DGraph); err != nil {
		log.Panicf("[MAIN] failed to open d-graph: %s", err.Error())
	}

	// rest
	if err := server.Run("localhost:17073"); err != nil {
		log.Panicf("[MAIN] failed to start rest server: %s", err.Error())
	}
}
