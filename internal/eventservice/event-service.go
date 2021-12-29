package eventservice

import (
	"context"
	"net/http"
	"time"

	"github.com/crazy-max/swarm-cronjob/internal/app"
	"github.com/rs/zerolog/log"
	"gopkg.in/macaron.v1"
)

type EventService struct {
	m       *macaron.Macaron
	srv     *http.Server
	port    string
	timeout string
}

func NewEventService(app *app.SwarmCronjob, port string, timeout string) *EventService {

	m := macaron.Classic()
	m.Map(app)

	m.Get("/event/:service/:key", handle)

	es := &EventService{
		m:       m,
		port:    port,
		timeout: timeout,
	}

	return es
}

func handle(app *app.SwarmCronjob, ctx *macaron.Context) (int, string) {
	serviceName := ctx.Params(":service")
	serviceKey := ctx.Params(":key")
	log.Info().Str("service", serviceName)

	tasks, err := app.Tasks(serviceName)
	if err != nil {
		return 500, err.Error()
	}
	app.RunJobByEvent(serviceName, serviceKey)

	time.Sleep(5 * time.Second)

	msg, err := app.WaitForEnd(serviceName, tasks)

	if err != nil {
		return 500, err.Error()
	} else {
		return 200, msg
	}
}

func (es *EventService) Run() {

	timeout, err := time.ParseDuration(es.timeout)
	if err != nil {
		log.Fatal().Msgf("Can not parse event timeout:%+s\n", err)
	}

	srv := &http.Server{
		Handler: es.m,
		Addr:    ":" + es.port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: timeout,
		ReadTimeout:  15 * time.Second,
	}

	log.Info().Msgf("Запуск на порту %s", es.port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("listen:%+s\n", err)
		}
	}()
}

func (es *EventService) Shutdown() {
	if es.srv != nil {
		es.srv.Shutdown(context.Background())
	}
}
