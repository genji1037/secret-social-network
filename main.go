package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/vrecan/death"
	_ "net/http/pprof"
	"os"
	"secret-social-network/app/config"
	"secret-social-network/app/dgraph"
	"secret-social-network/app/server"
	"secret-social-network/app/service"
	"secret-social-network/app/storage"
	"server-open-api/models/exchange"
	"sync/atomic"
	"syscall"
	"time"
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

	go func() {
		// 捕捉退出信号
		d := death.NewDeath(syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL,
			syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM)
		d.WaitForDeathWithFunc(func() {
			startAt := time.Now()
			timeout := 5 * time.Second
			for {
				if atomic.CompareAndSwapInt64(&service.InflightConsRelation, 0, -1) || time.Now().Sub(startAt) > timeout {
					break
				} else {
					log.Infof("Receive exit signal, waiting [%d] inflight cons relation ...", exchange.RewardingCnt)
					time.Sleep(100 * time.Millisecond)
				}
			}
		})
		log.Infof("ssns server stopped.")
		os.Exit(0)
	}()

	// rest
	if err := server.Run("localhost:17073"); err != nil {
		log.Panicf("[MAIN] failed to start rest server: %s", err.Error())
	}
}
