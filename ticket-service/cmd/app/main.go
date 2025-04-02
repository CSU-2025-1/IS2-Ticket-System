package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"ticket-service/cmd/migrator"
	"ticket-service/internal/httpengine"
	"ticket-service/internal/httpengine/handler"
	kafkaRepo "ticket-service/internal/repository/kafka"
	"ticket-service/internal/repository/postgres"
	"ticket-service/pkg/consul"
	"time"
)

const connection = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"

type Config struct {
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"postgres"`

	Kafka struct {
		Broker string `yaml:"broker"`
		Topic  string `yaml:"topic"`
	} `yaml:"kafka"`

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

	serviceUuuid, err := consul.Register("public-ticker-service", "ticket", 8080)
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
		"false",
	)

	err = migrator.MigratePostgres(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	ticketRepository := postgres.NewTicketRepository(db)

	var conn *kafka.Conn
	for {
		conn, err = kafka.DialLeader(context.Background(), "tcp", config.Kafka.Broker, config.Kafka.Topic, 0)
		if err != nil {
			log.Println("error connecting to kafka:", err)
			time.Sleep(1 * time.Second)
			continue
		} else {
			break
		}
	}

	ticketSaver := kafkaRepo.NewTicketSaver(conn)

	h := handler.New(ticketRepository, ticketSaver)

	r := httpengine.NewRouter(h)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
