package main

import (
  "log"
  "fmt"
  "errors"
  "strings"
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

  if commonMsg.ObjectType != "rules_pic" && !strings.Contains(commonMsg.EventType, "click") {
    log.Println("not a good request")
    return nil
  }

  playerID := commonMsg.QueryParams["player_id"]
  gameID := commonMsg.QueryParams["game_id"]
  column := commonMsg.QueryParams["column_id"]

  log.Println(gameID)
  validNumbers := map[string]int{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7}

  _, ok := validNumbers[column]

  if ok && gameID != "" {
    columnName := fmt.Sprintf("col_%s", column)
    queryTemplate := `
      INSERT INTO connect_four_games (game_id, %s, active_player)
      VALUES (%s, ARRAY[%s], %s)
      ON CONFLICT (game_id)
      DO UPDATE SET %s = array_append(connect_four_games.%s, %s), active_player = EXCLUDED.active_player
      WHERE NOT EXISTS (SELECT * FROM connect_four_games WHERE active_player = %s AND game_id = %s);
    `

    query := []string{
      fmt.Sprintf(
        queryTemplate,
        columnName,
        gameID,
        playerID,
        playerID,
        columnName,
        columnName,
        playerID,
        playerID,
        gameID,
      ),
    }

    fmt.Printf("%s", query)

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
