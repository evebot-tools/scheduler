package internal

import (
	"context"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/evebot-tools/queue"
	"github.com/rs/zerolog/log"
	"strconv"
)

func scheduleUniverseTypes() {
	log.Info().Msg("Scheduling universe types job")
	_, err := scheduler.Every(1).Day().At("11:10").Tag("universe_types").Do(universeTypesHandler)
	if err != nil {
		log.Err(err).Msg("Failed to schedule universe types")
	}
}

func universeTypesHandler() {
	log.Info().Msg("Publishing universe types")
	var typeInts []int32
	_, resp, err := eve.ESI.UniverseApi.GetUniverseTypes(context.Background(), nil)
	if err != nil {
		log.Err(err).Msg("Failed to get universe types")
		return
	}
	pages, err := strconv.Atoi(resp.Header.Get("x-pages"))
	if err != nil {
		log.Err(err).Msg("Failed to get x-pages")
		return
	}
	log.Debug().Any("pages", pages).Msg("got pages")

	for i := 1; i <= pages; i++ {
		t, _, err := eve.ESI.UniverseApi.GetUniverseTypes(context.Background(), &esi.GetUniverseTypesOpts{Page: optional.NewInt32(int32(i))})
		if err != nil {
			log.Err(err).Msg("Failed to get universe types")
			continue
		}
		typeInts = append(typeInts, t...)
	}

	for _, i := range typeInts {
		msg := queue.NewCronJobWithIntData("universe_types", true, 86400, int(i))
		err := ne.Publish(queue.TOPIC_JOB_UNIVERSE_TYPE, msg)
		if err != nil {
			log.Err(err).Msg("Failed to publish universe types")
		}
	}

}
