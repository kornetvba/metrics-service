package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	db2 "github.com/kornetvba/metrics-service/internal/config/db"
	"github.com/kornetvba/metrics-service/internal/config/logger"
	"github.com/kornetvba/metrics-service/internal/config/server"
	handlers "github.com/kornetvba/metrics-service/internal/handler"
	"github.com/kornetvba/metrics-service/internal/storage"
	"log"
	"net/http"
)

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server.ParseFlagServer()
	err := logger.Initialization()
	db := db2.Database{nil}
	_, err = db.New(server.DatabaseDSN)
	if err != nil {
		log.Print(err)
	}

	if err != nil {
		log.Print(err)
	}
	log.Print(server.FilePathStorage)
	err = run() // сервер
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Printf("serv is running %s", server.AddrServer.String())

	memStorage := storage.NewMemStorage()

	fileStore := storage.NewFileStorage(server.FilePathStorage, server.StorageInterval, server.Restore, memStorage)

	handler := handlers.NewMetricHandler(fileStore)
	//upload from a file

	err := fileStore.Load()
	if err != nil {
		log.Print(err)
	}

	// Канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем периодическое сохранение в отдельной горутине
	if fileStore.StorageIntervalSave != 0 {
		go func() {
			for {
				time.Sleep(time.Duration(server.StorageInterval) * time.Second)
				fmt.Println("save kd")
				fileStore.Save()
			}
		}()
	}

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(server.GzipMiddleware)
		r.Use(logger.LogMiddlewareGet)
		r.Get("/", handler.GetAllMetricsHTML)
		r.Get("/ping", handler.PingHandler)
		r.Post("/value/", handler.MetricGetJSON)
		r.Get("/value/{type_metric}/{name_metric}", handler.MetricGet)
	})

	r.Route("/update", func(r chi.Router) {
		r.Use(server.GzipMiddleware)
		r.Use(logger.LogMiddlewarePost)
		r.Use(fileStore.SaveToFile)
		r.Post("/", handler.MetricPostJSON)
		r.Post("/{type_metric}/{name_metric}/{value_metric}", handler.MetricPost)
	})

	// Создаем HTTP сервер с таймаутами
	srv := &http.Server{
		Addr:    server.AddrServer.String(),
		Handler: r,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on %s", server.AddrServer.String())

	// Ждем сигнал завершения
	<-done
	log.Println("Server is shutting down...")

	// Сохраняем данные перед выходом
	log.Println("Saving data to file...")
	err = fileStore.Save()
	if err != nil {
		log.Print("Save to file not success")
	}
	log.Println("Data saved successfully")

	// Даем серверу время завершить текущие запросы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server exited properly")
	return nil
}
