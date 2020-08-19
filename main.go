package main

import (
	"log"
	"net/http"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/stormdynamics/swerpctask/src/config"
	"github.com/stormdynamics/swerpctask/src/incrementor"
	"github.com/stormdynamics/swerpctask/src/rpchttpjson"
	"github.com/stormdynamics/swerpctask/src/rpcapi.v1"
	"github.com/stormdynamics/swerpctask/src/migration"
)

func main() {
	cfg := config.NewConfig()
	
	db, err := sql.Open("mysql", cfg.DbDSN())

	if err != nil {
		log.Fatalf("database open error: %s\n", err.Error())
	}

	defer db.Close()

	migration.Migration(db)
	
	inc := incrementor.NewIncrementor(db)
	api := rpcapi.NewApi(inc)
	mux := http.NewServeMux()
	
	// используется самый простой стандартый ServeMux
	// можно заменить на любой другой HTTP роутер (мультиплексор)
	mux.Handle(cfg.AppApiPath(), rpchttpjson.NewRpcHttpJson(api))

	srv := http.Server{
		Addr:    cfg.AppHostPort(),
		Handler: mux,
	}

	// создаем канал для получения сигналов от ОС
	done := make(chan os.Signal, 1)

	// подписываемся на сигналы SIGINT, SIGTERM
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// запуск HTTP сервера в отдельной go-рутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error %s\n", err)
		}
	}()

	log.Printf("http server started on %s\n", srv.Addr)

	// ожидаем поступление сигналов SIGINT, SIGTERM
	<-done

	log.Println("http server stopped")

	// завершаем работу, вводим таймер для таймаута завершения равным 5 секундам

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("http server shutdown failed: %s\n", err.Error())
	} else {
		log.Println("http server shutdown success")
	}
}