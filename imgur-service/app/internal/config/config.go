package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug       *bool `yaml:"is_debug" env:"IMG_IS_DEBUG" env-default:"false"  env-required:"true"`
	IsDevelopment *bool `yaml:"is_development" env:"IMG_IS_DEVELOPMENT" env-default:"false" env-required:"true"`
	Imgur         struct {
		AccessToken  string `yaml:"access_token" env:"IMG_ACCESS_TOKEN" env-required:"true"`
		ClientSecret string `yaml:"client_secret" env:"IMG_CLIENT_SECRET" env-required:"true"`
		ClientID     string `yaml:"client_id" env:"IMG_CLIENT_ID" env-required:"true"`
		URL          string `yaml:"url" env:"IMG_URL" env-required:"true"`
	} `yaml:"imgur"`
	RabbitMQ struct {
		Host     string `yaml:"host" env:"IMG_RABBIT_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"IMG_RABBIT_PORT" env-required:"true"`
		Username string `yaml:"username" env:"IMG_RABBIT_USERNAME" env-required:"true"`
		Password string `yaml:"password" env:"IMG_RABBIT_PASSWORD" env-required:"true"`
		Consumer struct {
			Queue              string `yaml:"queue" env:"IMG_RABBIT_CONSUMER_QUEUE" env-required:"true"`
			MessagesBufferSize int    `yaml:"messages_buffer_size" env:"IMG_RABBIT_CONSUMER_MBS" env-required:"true"`
		} `yaml:"consumer"`
		Producer struct {
			Queue string `yaml:"queue" env:"IMG_RABBIT_PRODUCER_QUEUE" env-required:"true"`
		} `yaml:"producer"`
	} `yaml:"rabbit_mq"`
	AppConfig AppConfig `yaml:"app" env-required:"true"`
}

type AppConfig struct {
	EventWorkers int    `yaml:"event_workers" env:"IMG_EVENT_WORKERS" env-default:"3" env-required:"true"`
	LogLevel     string `yaml:"log_level" env:"IMG_LOG_LEVEL" env-default:"error" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig(path string) *Config {
	once.Do(func() {
		log.Printf("read application config in path %s", path)

		instance = &Config{}

		if err := cleanenv.ReadConfig(path, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
