package main

import (
  "log"
	"net/http"

  "github.com/joho/godotenv"
  "github.com/julienschmidt/httprouter"

  auth "github.com/ianlofs/go_rest_api/authentication"
  "github.com/ianlofs/go_rest_api/database"
)

func main() {

  db, err := database.New()
  if err != nil {
    log.Fatalln(err)
  }
  db.SetMaxOpenConns(40)
  defer db.Close()

  auth.InitAuth()

  router := httprouter.New()
  router.POST("/login", auth.AuthUser(db))
  router.POST("/", auth.AuthRequest)

  log.Fatal(http.ListenAndServe(":8080", router))
}
