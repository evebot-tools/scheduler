package internal

import (
	"context"
	"github.com/antihax/goesi"
	"github.com/evebot-tools/utils"
	"github.com/go-co-op/gocron"
	"github.com/gomodule/redigo/redis"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	Jobs      = make(map[string]*gocron.Job)
	scheduler *gocron.Scheduler
	apiServer *http.Server
	eve       *goesi.APIClient
	nc        *nats.Conn
	ne        *nats.EncodedConn
)

func Init() {
	utils.InitMongoDBClient()
}

func Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	log.Info().Msg("Configuring scheduler")
	scheduler = gocron.NewScheduler(time.UTC)
	scheduler.WaitForScheduleAll()

	var err error
	var esiCache redis.Conn
	eve, esiCache = utils.InitEsiClient()

	nc, err = nats.Connect(utils.GetEnv("NATS_ADDR", "nats://localhost:4222"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to nats")
	}
	ne, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to nats with json encoder")
	}

	scheduleStatus()
	scheduleUniverseTypes()

	scheduler.StartAsync()

	go startApi()

	for {
		select {
		case <-quit:
			log.Info().Msg("SIGINT received, exiting")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = apiServer.Shutdown(ctx)
			if err != nil {
				log.Err(err).Msg("failed to shutdown webserver")
			}
			scheduler.Stop()
			err = esiCache.Close()
			if err != nil {
				log.Err(err).Msg("failed to close eve cache connection")
			}
			ne.Close()
			nc.Close()
		}
	}
}
