package main

import (
  cfg "github.com/movableink/go-go-common-go/config"
  "github.com/movableink/go-go-common-go/pg"
)

type ConnectFourConfig struct {
  Topic           string
  Channel         string
}

type BaseConfig struct {
  cfg.Config
  ConnectFour ConnectFourConfig
  PgConfig pg.Config
}

var Config = &BaseConfig{}

func (c *BaseConfig) Load() {
  c.MarshalKey("connectFour", &c.ConnectFour)
  c.MarshalKey("pgConfig", &c.PgConfig)
}
