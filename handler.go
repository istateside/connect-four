package main

import (
  "log"
  "fmt"
  "errors"
  "encoding/json"
  common "github.com/movableink/go-go-common-go"
  "github.com/movableink/go-go-common-go/pg"
  "github.com/bitly/go-nsq"
)

type Handler struct {
  pg *pg.Pg
}

func NewHandler() *Handler {
  return &Handler{
    pg: pg.New(&Config.PgConfig),
  }
}

func (h *Handler) Handle(msg *nsq.Message) error {
  commonMsg := &common.Message{}
  err := json.Unmarshal(msg.Body, commonMsg)

  if err != nil {
    log.Println("shit this broke")
    return err
  }

  playerID := commonMsg.QueryParams["player_id"]
  gameID := commonMsg.QueryParams["game_id"]
  column := commonMsg.QueryParams["column_idx"]

  validNumbers := map[string]int{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7}

  _, ok := validNumbers[column]

  if ok && gameID != "" {
    columnName := fmt.Sprintf("column_%s", column)
    queryTemplate := "update connect_four_games set %s = array_append(%s, '%s') where game_id = %s"
    query := []string{fmt.Sprintf(queryTemplate, columnName, columnName, playerID, gameID)}
    h.pg.Write(query)
  } else {
    err := errors.New("who knows")
    log.Println(err.Error())
    return err
  }

  return nil
}

func (h *Handler) Start() {}
func (h *Handler) SetConsumer(c *nsq.Consumer) {}
func (h *Handler) Shutdown() {
  h.pg.Shutdown()
}
