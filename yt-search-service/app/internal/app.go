package internal

import (
	"github.com/julienschmidt/httprouter"
	"github.com/serjic/yt-search/internal/config"
	"github.com/serjic/yt-search/internal/events"
	youtubeService "github.com/serjic/yt-search/internal/youtube"
	"github.com/serjic/yt-search/pkg/client/mq/rabbitmq"
	"github.com/serjic/yt-search/pkg/client/youtube"
	"github.com/serjic/yt-search/pkg/logger"
	"net/http"
)

type app struct {
	cfg            *config.Config
	logger         *logger.Logger
	httpServer     *http.Server
	youtubeService youtubeService.Service
	router         *httprouter.Router
}

type App interface {
	Run()
}

func NewApp(logger *logger.Logger, cfg *config.Config) (App, error) {
	logger.Println("router initializing")
	router := httprouter.New()

	ytClient := youtube.NewClient(cfg.Youtube.APIURL, cfg.Youtube.AccessToken, &http.Client{})
	yts := youtubeService.NewService(logger, ytClient)

	return &app{
		cfg:            cfg,
		logger:         logger,
		youtubeService: yts,
		router:         router,
	}, nil
}

func (a *app) Run() {
	a.startConsume()
}

func (a *app) startConsume() {
	a.logger.Info("start consuming")

	consumer, err := rabbitmq.NewRabbitMQConsumer(rabbitmq.ConsumerConfig{
		BaseConfig: rabbitmq.BaseConfig{
			Host:     a.cfg.RabbitMQ.Host,
			Port:     a.cfg.RabbitMQ.Port,
			Username: a.cfg.RabbitMQ.Username,
			Password: a.cfg.RabbitMQ.Password,
		},
		PrefetchCount: a.cfg.RabbitMQ.Consumer.MessagesBufferSize,
	})
	if err != nil {
		a.logger.Fatal(err)
	}
	producer, err := rabbitmq.NewRabbitMQProducer(rabbitmq.ProducerConfig{
		BaseConfig: rabbitmq.BaseConfig{
			Host:     a.cfg.RabbitMQ.Host,
			Port:     a.cfg.RabbitMQ.Port,
			Username: a.cfg.RabbitMQ.Username,
			Password: a.cfg.RabbitMQ.Password,
		},
	})
	if err != nil {
		a.logger.Fatal(err)
	}

	messages, err := consumer.Consume(a.cfg.RabbitMQ.Consumer.Queue)
	if err != nil {
		a.logger.Fatal(err)
	}

	client := http.Client{}
	ytClient := youtube.NewClient(a.cfg.Youtube.APIURL, a.cfg.Youtube.AccessToken, &client)
	ytService := youtubeService.NewService(a.logger, ytClient)

	for i := 0; i < a.cfg.AppConfig.EventWorkers; i++ {
		worker := events.NewWorker(i, consumer, a.cfg.RabbitMQ.Producer.Queue, producer, messages, a.logger, ytService)

		go worker.Process()
		a.logger.Infof("Event Worker #%d started", i)
	}
}
