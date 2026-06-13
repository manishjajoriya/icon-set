package config

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Http HttpConfig
	Icon IconConfig
}

type HttpConfig struct {
	Port string `env:"HTTP_PORT" env-required:"true"`
}

type IconConfig struct {
	IconJsonPath        string   `env:"ICON_JSON_PATH" env-required:"true"`
	OnlyMultiColor      bool     `env:"ONLY_MULTI_COLOR" env-required:"true"`
	DefaultSize         float32  `env:"DEFAULT_SIZE" env-required:"true" env-default:"24"`
	AllowedIconPackPath string   `env:"ALLOWED_ICON_PACK_PATH" env-default:""`
	AllowedIconPack     []string `json:"allowed_icon_pack"`
}

type allowedIconPack struct {
	AllowedIconPack []string `json:"allowed_icon_pack"`
}

func MustLoad() *Config {
	var configPath string
	var cfg *Config

	cfg = &Config{}

	err := cleanenv.ReadEnv(cfg)

	if err == nil {
		loadAllowedIconPack(cfg)
		return cfg
	}

	log.Info().Msg("failed to read env, falling back to CONFIG_PATH")

	configPath = strings.TrimSpace(os.Getenv("CONFIG_PATH"))

	if configPath == "" {
		flags := flag.String("config", ".env", "config file path")
		flag.Parse()
		configPath = *flags
	}

	if configPath == "" {
		log.Info().Msg("config file path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal().Msg("config file does not exist")
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatal().Err(err).Msg("error reading config")
	}

	log.Info().
		Str("path", cfg.Icon.AllowedIconPackPath).
		Msg("loading allowed icon packs")
	loadAllowedIconPack(cfg)

	return cfg
}

func loadAllowedIconPack(cfg *Config) {
	if cfg.Icon.AllowedIconPackPath == "" {
		return
	}

	data, err := os.ReadFile(cfg.Icon.AllowedIconPackPath)
	if err != nil {
		log.Fatal().Err(err).Msg("error reading allowed icon pack")
	}

	var al allowedIconPack
	if err := json.Unmarshal(data, &al); err != nil {
		log.Fatal().Err(err).Msg("error unmarshaling allowed icon pack")
	}

	cfg.Icon.AllowedIconPack = al.AllowedIconPack
}
