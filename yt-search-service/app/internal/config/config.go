package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug       *bool `yaml:"is_debug" env:"DATA-YTS-IsDebug" env-default:"false"  env-required:"true"`
	IsDevelopment *bool `yaml:"is_development" env:"DATA-BOT-IsDevelopment" env-default:"false" env-required:"true"`
	Listen        struct {
		Type   string `yaml:"type" env:"DATA-YTS-ListenType" env-default:"port"`
		BindIP string `yaml:"bind_ip" env:"DATA-YTS-BindIP" env-default:"localhost"`
		Port   string `yaml:"port" env:"DATA-YTS-Port" env-default:"8080"`
	} `yaml:"listen" env-required:"true"`
	Youtube struct {
		APIURL          string `yaml:"api_url" env:"DATA-YTS-YoutubeAPIURL" env-required:"true"`
		RefreshTokenURL string `yaml:"refresh_token_URL" env:"DATA-YTS-RefreshTokenURL" env-required:"true"`
		APIKey          string `yaml:"api_key" env:"DATA-YTS-APIKey" env-required:"true"`
		ClientID        string `yaml:"client_id" env:"DATA-YTS-ClientID" env-required:"true"`
		ClientSecret    string `yaml:"client_secret" env:"DATA-YTS-ClientSecret" env-required:"true"`
		AccessToken     string `yaml:"access_token" env:"DATA-YTS-YoutubeAccessToken" env-required:"true"`
		RefreshToken    string `yaml:"refresh_token" env:"DATA-YTS-RefreshToken" env-required:"true"`
		AuthSuccessUri  string `yaml:"auth_success_uri" env:"DATA-YTS-AuthSuccessUri" env-required:"true"`
		AccountsUri     string `yaml:"accounts_uri" env:"DATA-YTS-AccountsUri" env-required:"true"`
	} `yaml:"youtube" env-required:"true"`
	RabbitMQ struct {
		Host     string `yaml:"host" env:"DATA-YTS-RabbitHost" env-required:"true"`
		Port     string `yaml:"port" env:"DATA-YTS-RabbitPort" env-required:"true"`
		Username string `yaml:"username" env:"DATA-YTS-RabbitUsername" env-required:"true"`
		Password string `yaml:"password" env:"DATA-YTS-RabbitPassword" env-required:"true"`
		Consumer struct {
			Queue              string `yaml:"queue" env:"DATA-YTS-RabbitConsumerQueue" env-required:"true"`
			MessagesBufferSize int    `yaml:"messages_buffer_size" env:"DATA-YTS-RabbitConsumerMBS" env-required:"true"`
		} `yaml:"consumer" env-required:"true"`
		Producer struct {
			Queue string `yaml:"queue" env:"DATA-YTS-RabbitProducerQueue" env-required:"true"`
		} `yaml:"producer" env-required:"true"`
	} `yaml:"rabbit_mq" env-required:"true"`
	AppConfig AppConfig `yaml:"app" env-required:"true"`
}

type AppConfig struct {
	EventWorkers int    `yaml:"event_workers" env:"DATA-YTS-EventWorks" env-default:"3" env-required:"true"`
	LogLevel     string `yaml:"log_level" env:"DATA-YTS-LogLevel" env-default:"error" env-required:"true"`
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
