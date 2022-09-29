package bvbot

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

func NewConfig() (cfg Config) {
	cfg.Courts = NewCourtsConfig()
	return
}

type Config struct {
	Courts CourtsConfig
}

func (conf Config) Value() (driver.Value, error) {
	return json.Marshal(conf)
}

func (conf *Config) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &conf)
}

type CourtsConfig struct {
	Max        int `json:"max"`
	MaxPlayers int `json:"max_players"`
	MinPlayers int `json:"min_players"`
}

func NewCourtsConfig() CourtsConfig {
	return CourtsConfig{Max: 4, MaxPlayers: 12, MinPlayers: 4}
}
