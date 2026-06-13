package util

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"

	"github.com/manishjajoriya/icon-set/internal/config"
	"github.com/rs/zerolog/log"
)

var colorIconPack = []string{
	"catppuccin.json", "cif.json", "circle-flags.json", "cryptocurrency-color.json", "devicon-original.json",
	"devicon.json", "emojione-v1.json", "emojione.json", "flag.json", "flagpack.json", "flat-color-icons.json",
	"flat-ui.json", "fluent-color.json", "fluent-emoji-flat.json", "fluent-emoji.json", "fxemoji.json", "gcp.json",
	"glyphs-poly.json", "icon-park.json", "logos.json", "marketeq.json", "material-icon-theme.json", "meteocons.json",
	"noto-v1.json", "noto.json", "openmoji.json", "skill-icons.json", "streamline-color.json",
	"streamline-cyber-color.json", "streamline-emojis.json", "streamline-flex-color.json", "streamline-freehand-color.json",
	"streamline-kameleon-color.json", "streamline-plump-color.json", "streamline-sharp-color.json",
	"streamline-stickies-color.json", "streamline-ultimate-color.json", "token-branded.json", "twemoji.json",
	"unjs.json", "vscode-icons.json",
}

type Icon struct {
	Key    string
	Name   string
	Prefix string
	Body   string
	Width  float32
	Height float32
}

type InfoStruct struct {
	Palette bool `json:"palette"`
}

type iconifyFile struct {
	Prefix string `json:"prefix"`
	Icons  map[string]struct {
		Body   string  `json:"body"`
		Width  float32 `json:"width"`
		Height float32 `json:"height"`
	} `json:"icons"`
	Height float32    `json:"height"`
	Info   InfoStruct `json:"info"`
}

func LoadIcons(cfg *config.Config) (map[string]Icon, error) {
	files, err := filepath.Glob(cfg.Icon.IconJsonPath)
	if err != nil {
		return nil, err
	}

	icons := make(map[string]Icon)

	for _, file := range files {
		if cfg.Icon.OnlyMultiColor == true {
			found := slices.Contains(colorIconPack, filepath.Base(file))
			if found == false {
				continue
			}
		}

		if len(cfg.Icon.AllowedIconPack) != 0 {
			if found := slices.Contains(cfg.Icon.AllowedIconPack, filepath.Base(file)); found == false {
				continue
			}
		}

		log.Info().Str("file", filepath.Base(file)).Msg("loading icon")

		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		var collection iconifyFile
		if err := json.Unmarshal(data, &collection); err != nil {
			return nil, err
		}

		for name, icon := range collection.Icons {
			key := collection.Prefix + ":" + name

			var height float32
			if icon.Height != 0 {
				height = icon.Height
			} else if collection.Height != 0 {
				height = collection.Height
			} else {
				height = -1
			}

			var width float32
			if icon.Width != 0 {
				width = icon.Width
			} else {
				width = -1
			}

			icons[key] = Icon{
				Key:    key,
				Name:   name,
				Prefix: collection.Prefix,
				Body:   icon.Body,
				Width:  width,
				Height: height,
			}
		}
	}

	log.Info().Int("length", len(icons)).Msg("loaded icons")

	return icons, nil
}

func LoadMultiColorIcons() (map[string]Icon, error) {
	files, err := filepath.Glob("./internal/assets/json/*.json")
	if err != nil {
		return nil, err
	}

	colors := make([]string, 0)

	icons := make(map[string]Icon)

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		var collection iconifyFile
		if err := json.Unmarshal(data, &collection); err != nil {
			return nil, err
		}

		if collection.Info.Palette {
			colors = append(colors, filepath.Base(file))
		}

	}

	log.Info().Int("length", len(icons)).Msg("loaded multicolor icons")
	log.Info().Int("icon pack length", len(colors)).Msg("loaded multicolor icons")
	log.Info().Any("icon pack", colors).Msg("loaded multicolor icons")

	return icons, nil
}
