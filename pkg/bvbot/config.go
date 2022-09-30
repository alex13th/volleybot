package bvbot

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func NewConfig() (cfg Config) {
	cfg.Courts = NewCourtsConfig()
	cfg.Price = NewPriceConfig()
	return
}

type Config struct {
	Courts CourtsConfig
	Price  PriceConfig
}

func (conf Config) Value() (driver.Value, error) {
	return json.Marshal(conf)
}

func (conf *Config) Scan(value interface{}) (err error) {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	if err = json.Unmarshal(b, &conf); err != nil {
		return
	}
	estr := []string{}
	if (conf.Courts == CourtsConfig{}) {
		conf.Courts = NewCourtsConfig()
		estr = append(estr, "empty CourtsConfig")
	}
	if (conf.Price == PriceConfig{}) {
		conf.Price = NewPriceConfig()
		estr = append(estr, "empty PriceConfig")
	}
	if len(estr) > 0 {
		err = errors.New(strings.Join(estr, ", "))
	}
	return
}

type CourtsConfig struct {
	Max        int `json:"max"`
	MaxPlayers int `json:"max_players"`
	MinPlayers int `json:"min_players"`
}

func NewCourtsConfig() CourtsConfig {
	return CourtsConfig{Max: 4, MaxPlayers: 12, MinPlayers: 4}
}

type PriceConfig struct {
	Min  int
	Max  int
	Step int
}

func NewPriceConfig() PriceConfig {
	return PriceConfig{Min: 0, Max: 2000, Step: 100}
}

type ConfigTelegramView struct {
	Config
	ParseMode string
}

func NewConfigTelegramViewRu(cfg Config) ConfigTelegramView {
	return ConfigTelegramView{Config: cfg, ParseMode: "Markdown"}
}

func (tgv ConfigTelegramView) GetText() (text string) {
	text = NewConfigCourtsTelegramViewRu(tgv.Config.Courts).GetText()
	text += "\n\n"
	text += NewConfigPriceTelegramViewRu(tgv.Config.Price).GetText()
	return
}

type ConfigCourtsTelegramView struct {
	CourtsConfig
	Resources ConfigCourtsResources
	ParseMode string
}

func NewConfigCourtsTelegramViewRu(cfg CourtsConfig) ConfigCourtsTelegramView {
	return ConfigCourtsTelegramView{
		CourtsConfig: cfg,
		Resources:    NewConfigCourtsResourcesRu(),
		ParseMode:    "Markdown",
	}
}

func (tgv ConfigCourtsTelegramView) GetText() (text string) {
	text = "⚙️*Настройки кортов:*"
	text += fmt.Sprintf("\n*%s*: %v", tgv.Resources.Max, tgv.CourtsConfig.Max)
	text += fmt.Sprintf("\n*%s*: %v", tgv.Resources.MinPlayers, tgv.CourtsConfig.MinPlayers)
	text += fmt.Sprintf("\n*%s*: %v", tgv.Resources.MaxPlayers, tgv.CourtsConfig.MaxPlayers)
	return
}

type ConfigPriceTelegramView struct {
	PriceConfig
	Resources ConfigPriceResources
	ParseMode string
}

func NewConfigPriceTelegramViewRu(cfg PriceConfig) ConfigPriceTelegramView {
	return ConfigPriceTelegramView{
		PriceConfig: cfg,
		Resources:   NewConfigPriceResourcesRu(),
		ParseMode:   "Markdown",
	}
}

func (tgv ConfigPriceTelegramView) GetText() (text string) {
	text = "⚙️*Настройки цены:*"
	text += fmt.Sprintf("\n*%s*: %v", tgv.Resources.Min, tgv.PriceConfig.Min)
	text += fmt.Sprintf("\n*%s*: %v", tgv.Resources.Max, tgv.PriceConfig.Max)
	text += fmt.Sprintf("\n*%s*: %v", tgv.Resources.Step, tgv.PriceConfig.Step)
	return
}
