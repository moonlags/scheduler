package main

import (
	"fmt"
	"log"
)

func main() {
	cfg, err := newConfig()
	if err != nil {
		panic(fmt.Sprintf("create config object: %v", err))
	}

	db, err := openDB(cfg.dbUri)
	if err != nil {
		panic(fmt.Sprintf("create db object: %v", err))
	}

	server, err := newServer(cfg, db)
	if err != nil {
		panic(fmt.Sprintf("create server object: %v", err))
	}

	log.Fatal(server.run(":" + cfg.port))
}
