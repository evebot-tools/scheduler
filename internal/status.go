package internal

import (
	"github.com/evebot-tools/queue"
	"github.com/rs/zerolog/log"
)

func scheduleStatus() {
	log.Info().Msg("Scheduling status job")
	_, err := scheduler.Every(30).Seconds().Tag("status_status").Do(statusHandler)
	if err != nil {
		log.Err(err).Msg("Failed to schedule status status")
	}
}

func statusHandler() {
	log.Info().Msg("Publishing status status")
	msg := queue.NewCronJob("status", true, 30)
	log.Debug().Any("msg", msg).Msg("generated status msg")
	err := ne.Publish(queue.TOPIC_CRONJOB_STATUS, msg)
	log.Debug().Err(err).Msg("published status msg")
	if err != nil {
		log.Err(err).Msg("Failed to publish status status")
	}
}
