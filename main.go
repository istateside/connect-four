package main

import (
  common "github.com/movableink/go-go-common-go"
  "github.com/movableink/go-go-common-go/cli"
)

func main() {
  cli := cli.New("heres my description!!!", Config)
  cli.Execute(func() {
    common.NewWorker(
      Config.ConnectFour.Topic,
      Config.ConnectFour.Channel,
      NewHandler(),
      Config,
    )
  }).Run(1)
}
