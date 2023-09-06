package main

import (
	"TESTShop"
	"TESTShop/pkg/broker"
	"TESTShop/pkg/cache"
	"TESTShop/pkg/handler"
	"TESTShop/pkg/repository"
	"TESTShop/pkg/service"
	"context"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err)
	}

	cfg := broker.ConfigNATS{
		ClusterID:        viper.GetString("nats.clusterId"),
		ClientIdConsumer: viper.GetString("nats.clientIdConsumer"),
		ClientIdProducer: viper.GetString("nats.clientIdProducer"),
		ChannelName:      viper.GetString("nats.channel"),
		NatsUrl:          viper.GetString("nats.natsURL"),
		DurableName:      viper.GetString("nats.durableName"),
		Subject:          viper.GetString("nats.subject"),
	}

	db, err := repository.NewPostgresDB(repository.ConfigDB{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	sc, err := broker.CreateProducer(cfg.ClientIdProducer, cfg.ClusterID, cfg.NatsUrl)
	if err != nil {
		logrus.Fatalf("failed to initialize producer %s", err)
	}

	cacheInMemory := make(map[string]TESTShop.OrderResponse)
	cacheMutex := sync.Mutex{}

	repos := repository.NewRepository(db)
	brokers := broker.NewBroker(sc)
	caches := cache.NewCache(cacheInMemory, &cacheMutex)
	services := service.NewOrderService(repos, brokers, caches)
	handlers := handler.NewHandler(services)

	//заполнить конфиг

	go func() {
		err := services.StartProcessMessages(cfg.ChannelName)
		logrus.Print("start process messages")
		if err != nil {
			logrus.Errorf("error while processing messages %s", err)
		}
	}()

	err = broker.PublishMessage(sc, viper.GetString("nats.channel"), 2)
	if err != nil {
		logrus.Fatalf("error while sending messages %s", err)
	}

	services.InitCache()

	srv := new(TESTShop.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TestShopApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TestShopApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

}
