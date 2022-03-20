package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tele "gopkg.in/telebot.v3"

	"github.com/serj1c/tbot/internal/config"
	"github.com/serj1c/tbot/internal/events"
	"github.com/serj1c/tbot/internal/service"
	"github.com/serj1c/tbot/pkg/client/imgur"
	"github.com/serj1c/tbot/pkg/client/mq"
	"github.com/serj1c/tbot/pkg/client/mq/rabbitmq"
	"github.com/serj1c/tbot/pkg/logger"
)

type app struct {
	cfg          *config.Config
	logger       *logger.Logger
	httpServer   *http.Server
	imgurService service.ImgurService
	bot          *tele.Bot
	producer     mq.Producer
}

type App interface {
	Run()
}

func NewApp(logger *logger.Logger, cfg *config.Config) (App, error) {

	client := http.Client{}
	imgurClient := imgur.NewClient(cfg.Imgur.URL, cfg.Imgur.AccessToken, cfg.Imgur.ClientID, &client)
	imgurService := service.NewImgurService(logger, imgurClient)

	return &app{
		cfg:          cfg,
		logger:       logger,
		imgurService: imgurService,
	}, nil
}

func (a *app) Run() {
	a.startBot()
	a.startConsume()
	// TODO fixMe
	a.bot.Start()
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

	for i := 0; i < a.cfg.AppConfig.EventWorkers; i++ {
		worker := events.NewWorker(i, consumer, a.bot, producer, messages, a.logger)

		go worker.Process()
		a.logger.Infof("Event Worker #%d started", i)
	}

	a.producer = producer
}

func (a *app) startBot() {
	pref := tele.Settings{
		Token:   a.cfg.Telegram.Token,
		Poller:  &tele.LongPoller{Timeout: 60 * time.Second},
		Verbose: false,
		OnError: a.OnBotError,
	}
	var botErr error
	a.bot, botErr = tele.NewBot(pref)
	if botErr != nil {
		a.logger.Fatal(botErr)
		return
	}

	a.bot.Handle("/help", func(c tele.Context) error {
		return c.Send(fmt.Sprintf("/yt - search youtube for a track by its name"))
	})

	a.bot.Handle("/yt", func(c tele.Context) error {
		trackName := c.Message().Payload

		request := events.SearchTrackRequest{
			RequestID: fmt.Sprintf("%d", c.Sender().ID),
			Name:      trackName,
		}

		marshal, _ := json.Marshal(request)

		err := a.producer.Publish(a.cfg.RabbitMQ.Producer.Queue, marshal)
		if err != nil {
			return c.Send(fmt.Sprintf("error: %s", err.Error()))
		}

		return c.Send(fmt.Sprintf("Заявка принята"))
	})

	a.bot.Handle(tele.OnPhoto, func(c tele.Context) error {
		// Photos only.
		photo := c.Message().Photo
		file, err := a.bot.File(&photo.File)
		if err != nil {
			return c.Send("failed to download an image")
		}
		defer file.Close()
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(file)
		if err != nil {
			return c.Send("failed to download an image")
		}

		if buf.Len() > 10_485_760 {
			return c.Send("Лимит 10МБ")
		}

		image, err := a.imgurService.ShareImage(context.Background(), buf.Bytes())
		if err != nil {
			return c.Send("failed to share the image")
		}

		return c.Send(image)
	})
}

func (a *app) OnBotError(err error, ctx tele.Context) {
	a.logger.Error(err)
}
