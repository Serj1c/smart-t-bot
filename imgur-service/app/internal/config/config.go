package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug       *bool `yaml:"is_debug" env:"DATA-BOT-IsDebug" env-default:"false"  env-required:"true"`
	IsDevelopment *bool `yaml:"is_development" env:"DATA-BOT-IsDevelopment" env-default:"false" env-required:"true"`
	Listen        struct {
		Type   string `yaml:"type" env:"DATA-BOT-ListenType" env-default:"port"`
		BindIP string `yaml:"bind_ip" env:"DATA-BOT-BindIP" env-default:"localhost"`
		Port   string `yaml:"port" env:"DATA-BOT-Port" env-default:"8080"`
	} `yaml:"listen" env-required:"true"`
	Telegram struct {
		Token string `yaml:"token" env:"DATA-BOT-TelegramToken" env-required:"true"`
	}
	Imgur struct {
		AccessToken  string `yaml:"access_token" env:"DATA-BOT-ImgurAccessToken" env-required:"true"`
		ClientSecret string `yaml:"client_secret" env:"DATA-BOT-ImgurClientSecret" env-required:"true"`
		ClientID     string `yaml:"client_id" env:"DATA-BOT-ImgurClientID" env-required:"true"`
		URL          string `yaml:"url" env:"DATA-BOT-ImgurURL" env-required:"true"`
	} `yaml:"imgur"`
	RabbitMQ struct {
		Host     string `yaml:"host" env:"DATA-BOT-RabbitHost" env-required:"true"`
		Port     string `yaml:"port" env:"DATA-BOT-RabbitPort" env-required:"true"`
		Username string `yaml:"username" env:"DATA-BOT-RabbitUsername" env-required:"true"`
		Password string `yaml:"password" env:"DATA-BOT-RabbitPassword" env-required:"true"`
		Consumer struct {
			Queue              string `yaml:"queue" env:"DATA-BOT-RabbitConsumerQueue" env-required:"true"`
			MessagesBufferSize int    `yaml:"messages_buffer_size" env:"DATA-BOT-RabbitConsumerMBS" env-required:"true"`
		} `yaml:"consumer" env-required:"true"`
		Producer struct {
			Queue string `yaml:"queue" env:"DATA-BOT-RabbitProducerQueue" env-required:"true"`
		} `yaml:"producer" env-required:"true"`
	} `yaml:"rabbit_mq" env-required:"true"`
	AppConfig AppConfig `yaml:"app" env-required:"true"`
}

type AppConfig struct {
	EventWorkers int    `yaml:"event_workers" env:"DATA-BOT-EventWorks" env-default:"3" env-required:"true"`
	LogLevel     string `yaml:"log_level" env:"DATA-BOT-LogLevel" env-default:"error" env-required:"true"`
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
