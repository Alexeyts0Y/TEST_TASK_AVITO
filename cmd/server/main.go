package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/config"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/handler"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/model"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/repository"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDatabase(dsn string) *gorm.DB {
	var db *gorm.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Подключение к БД прошло успешно.")
			break
		}
		log.Printf("Ошибка подключения к БД (попытка %d/5). Повторное подключение через 2 секунды.", i+1)
		time.Sleep(2 * time.Second)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Team{}, &model.PullRequest{}); err != nil {
		log.Fatalf("Не удалось выполнить миграции: %v", err)
	}
	log.Printf("Миграции применились")

	if err != nil {
		log.Fatal("Не удалось подключиться к БД")
	}

	return db
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации БД: %v", err)
	}

	db := setupDatabase(cfg.DatabaseUrl)
	repository := repository.NewPostgresRepository(db)
	serviceHandler := handler.NewServer(repository)

	r := gin.Default()
	strictHandler := api.NewStrictHandler(serviceHandler, nil)

	api.RegisterHandlers(r, strictHandler)

	address := ":" + cfg.Port
	log.Printf("Сервер запускается на порту %v", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
