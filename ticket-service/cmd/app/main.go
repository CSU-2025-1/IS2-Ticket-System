package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"ticket-service/cmd/migrator"
	"ticket-service/internal/httpengine"
	"ticket-service/internal/httpengine/handler"
	"ticket-service/internal/repository/postgres"
	rabbitRepo "ticket-service/internal/repository/rabbitmq"
	"ticket-service/pkg/consul"
	"ticket-service/pkg/rabbitmq"
)

const connection = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"

type Config struct {
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"postgres"`

	Rabbit rabbitmq.Config `yaml:"rabbit"`

	Consul consul.Config `yaml:"consul"`
}

func loadConfig() (Config, error) {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	consul := consul.New(config.Consul)
	err = consul.Configure()
	if err != nil {
		log.Fatal(err)
	}

	serviceUuuid, err := consul.Register("public-ticket-service", "ticket", 8080)
	if err != nil {
		log.Fatal(err)
	}
	defer consul.Deregister(serviceUuuid)

	connStr := fmt.Sprintf(connection,
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.Database,
	)

	err = migrator.MigratePostgres(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	dbCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal(err)
	}

	dbCfg.MaxConns = 3

	db, err := pgxpool.NewWithConfig(context.Background(), dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	ticketRepository := postgres.NewTicketRepository(db)

	writer, err := rabbitmq.CreateJsonWriter[rabbitRepo.Ticket](config.Rabbit, rabbitmq.DefaultPublishOption)
	if err != nil {
		log.Fatal(err)
	}

	ticketSaver, err := rabbitRepo.NewRabbitMQ(writer)
	if err != nil {
		log.Fatal(err.Error())
	}

	h := handler.New(ticketRepository, ticketSaver)

	r := httpengine.NewRouter(h)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
