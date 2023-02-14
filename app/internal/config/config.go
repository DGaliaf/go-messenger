package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	HTTP struct {
		IP   string `yaml:"ip" env:"HTTP_IP" env-default:"localhost"`
		Port int    `yaml:"port" env:"HTTP_PORT" env-default:"9000"`
	} `yaml:"http"`
	PostgreSQL struct {
		Username string `yaml:"username" env:"POSTGRES_USER" env-required:"true" env-default:"postgres"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true" env-default:"postgres"`
		Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true" env-default:"0.0.0.0"`
		Port     string `yaml:"port" env:"POSTGRES_PORT" env-required:"true" env-default:"5432"`
		Database string `yaml:"database" env:"POSTGRES_DB" env-required:"true" env-default:"postgres"`
	} `yaml:"postgresql"`
}

const (
	configPath = "./configs/config.local.yml"
)

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Print("config init")

		instance = &Config{}

		if err := cleanenv.ReadConfig(configPath, instance); err != nil {
			log.Println("Cant`t read environment variables from neither .yaml nor .env")
			log.Println(err)

			err := cleanenv.ReadEnv(instance)
			if err != nil {
				help, _ := cleanenv.GetDescription(instance, nil)
				log.Println(help)
				log.Fatalln(err)
			}
		}
	})
	return instance
}
