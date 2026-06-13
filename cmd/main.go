package main

import (
	"time"

	"github.com/manishjajoriya/icon-set/internal/config"
	"github.com/manishjajoriya/icon-set/internal/util"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg := config.MustLoad()

	start := time.Now()
	icons, err := util.LoadIcons(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load icons")
	}

	timeTakeToLoad := float64(time.Since(start).Microseconds()) / 1000.0
	log.Info().Msgf("time take to load %v ms", timeTakeToLoad)

	app := app{
		icons: icons,
		cfg:   cfg,
	}
	app.Run(app.Init())
}
