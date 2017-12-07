package main

import (
  "log"
  "github.com/gorilla/mux"
  "database/sql"
  "net/http"
  "strconv"
  "encoding/json"
  _ "github.com/lib/pq"
)

type GameData struct {
  GameID int `json:"game_id"`
  ActivePlayer int `json:"active_player"`
  Col1 sql.NullString `json:"col_1"`
  Col2 sql.NullString `json:"col_2"`
  Col3 sql.NullString `json:"col_3"`
  Col4 sql.NullString `json:"col_4"`
  Col5 sql.NullString `json:"col_5"`
  Col6 sql.NullString `json:"col_6"`
  Col7 sql.NullString `json:"col_7"`
}

func main() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", getBoardState)
  log.Println("listening")
  http.ListenAndServe(":1990", router)
}

func getBoardState(w http.ResponseWriter, r *http.Request) {
  dbinfo := "postgres://kfleischman@localhost/tyson_dev?sslmode=disable"

  db, err := sql.Open("postgres", dbinfo)

  log.Println(dbinfo)

  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte(err.Error()))
    return
  }

  defer db.Close()

  queryValues := r.URL.Query()
  gameID, err := strconv.Atoi(queryValues.Get("game_id"))

  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("Game not found"))
    log.Println(err.Error())
    return
  }

  query := "SELECT * FROM connect_four_games WHERE game_id =$1"
  gameData := GameData{}
  err = db.QueryRow(query, gameID).Scan(
    &gameData.GameID,
    &gameData.ActivePlayer,
    &gameData.Col1,
    &gameData.Col2,
    &gameData.Col3,
    &gameData.Col4,
    &gameData.Col5,
    &gameData.Col6,
    &gameData.Col7,
  )

  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte(err.Error()))
    log.Printf("%s", queryValues)
    log.Println("Row for game not found")
    return
  }

  gameJson, err := json.Marshal(gameData)

  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte(err.Error()))
    return
  }

  w.Write(gameJson)
  return
}
