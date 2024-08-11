package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"test_task_BackDev/internal/handler"
	"test_task_BackDev/internal/repository"
	"test_task_BackDev/internal/server"
	"test_task_BackDev/internal/service"
	"test_task_BackDev/pkg/auth"
	"test_task_BackDev/pkg/database"
	"test_task_BackDev/pkg/email"
	"time"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "config"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Error(err)
		return
	}

	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
		return
	}

	db, err := database.NewPostgresDB(database.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),

		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	tokenManager, err := auth.NewManager(os.Getenv("SIGNING_KEY"))
	if err != nil {
		logrus.Error(err)
		return
	}

	emailSender, err := email.NewSMTPSender(
		viper.GetString("smtp.from"),
		os.Getenv("EMAIL_PASSWORD"),
		viper.GetString("smtp.host"),
		viper.GetInt("smtp.port"))

	repos := repository.NewRepository(db)
	services := service.NewService(service.Deps{
		Repos:           repos,
		TokenManager:    tokenManager,
		EmailSender:     emailSender,
		AccessTokenTLL:  viper.GetDuration("auth.accessTokenTLL"),
		RefreshTokenTLL: viper.GetDuration("auth.refreshTokenTLL"),
	})
	handlers := handler.NewHandler(services, tokenManager)

	srv := server.NewServer(viper.GetString("http.port"), handlers.InitRoutes())

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Run(); err != nil {
			logrus.Errorf("error while running http server: %s\n", err.Error())
		}
	}()

	logrus.Info("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdown()

		if err := srv.Stop(ctx); err != nil {
			logrus.Errorf("failed to stop server: %v", err.Error())
		}

		if err := db.Close(); err != nil {
			logrus.Errorf("failed to close database: %v", err.Error())
		}

		logrus.Info("Graceful shutdown complete")
	}()

	wg.Wait()
	logrus.Info("server stopped")
}

func initConfig() error {
	viper.AddConfigPath(CONFIG_DIR)
	viper.SetConfigName(CONFIG_FILE)
	return viper.ReadInConfig()
}
