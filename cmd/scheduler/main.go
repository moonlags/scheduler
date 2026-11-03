package main

import (
	"log"

	"github.com/moonlags/ikeaScheduler/internal/config"
	"github.com/moonlags/ikeaScheduler/internal/store"
	"github.com/moonlags/ikeaScheduler/internal/web"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	st, err := store.Open(cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer st.Close()

	if err := st.EnsureAdmin(cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Fatalf("bootstrap admin: %v", err)
	}

	srv, err := web.New(st)
	if err != nil {
		log.Fatalf("create server: %v", err)
	}

	if err := srv.Run(":" + cfg.Port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
