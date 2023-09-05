package main

import (
	"TESTShop"
	"TESTShop/pkg/broker"
	"TESTShop/pkg/handler"
	"TESTShop/pkg/repository"
	"TESTShop/pkg/service"
	"context"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err)
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

	sc, err := broker.CreateProducer(broker.ConfigNATSProducer{
		ClusterID: viper.GetString("natsProd.clusterId"),
		ClientID:  viper.GetString("natsProd.clientId"),
		NatsUrl:   viper.GetString("natsProd.natsURL"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize producer %s", err)
	}

	broker.SendMessage(sc, viper.GetString("natsProd.subject"), 5)

	err = broker.ReadFromChannel(broker.ConfigNATSConsumer{
		ClusterID:   viper.GetString("natsCons.clusterId"),
		ClientID:    viper.GetString("natsCons.clientId"),
		ChannelName: viper.GetString("natsCons.channel"),
		NatsUrl:     viper.GetString("natsCons.natsURL"),
		DurableName: viper.GetString("natsCons.durableName"),
	})

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

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

	logrus.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

}
