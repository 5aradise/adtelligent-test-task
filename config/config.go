package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Env       string    `yaml:"env"`
		DB        DB        `yaml:"db"`
		Server    Server    `yaml:"server"`
		Stitching Stitching `yaml:"stitching"`
	}

	DB struct {
		URL               string        `yaml:"url"`
		UpdateCacheDelay  time.Duration `yaml:"updateCacheDelay"`
		BigRequestTimeout time.Duration `yaml:"bigRequestTimeout"`
	}

	Server struct {
		AuctionPort   string        `yaml:"auctionPort"`
		StitchingPort string        `yaml:"stitchingPort"`
		Timeout       time.Duration `yaml:"timeout"`
		IdleTimeout   time.Duration `yaml:"idleTimeout"`
	}

	Stitching struct {
		AuctionUrl     string        `yaml:"auctionUrl"`
		RequestTimeout time.Duration `yaml:"requestTimeout"`
	}
)

func New(path string) (Config, error) {
	var c Config

	rawConfig, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("cannot load config file: %w", err)
	}

	err = yaml.Unmarshal(rawConfig, &c)
	if err != nil {
		return Config{}, fmt.Errorf("cannot unmarshal config file: %w", err)
	}

	return c, nil
}
