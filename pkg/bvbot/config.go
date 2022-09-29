package bvbot

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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
