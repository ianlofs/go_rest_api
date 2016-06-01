package main

import (
  "log"
	"net/http"

  "github.com/joho/godotenv"
  "github.com/julienschmidt/httprouter"

  "github.com/ianlofs/go_rest_api/constants"
  "github.com/ianlofs/go_rest_api/controllers"
  "github.com/ianlofs/go_rest_api/database"
  "github.com/ianlofs/go_rest_api/middleware"
)

func main() {
  godotenv.Load()
  // initialize private/public keys, db connection info, etc
  constants.InitConstants()

  db, err := database.New()
  if err != nil {
    log.Fatalln(err)
  }
  db.SetMaxOpenConns(40)
  defer db.Close()

  r := httprouter.New()
  r.POST("/register", controllers.Register(db))
  r.POST("/login", controllers.Login(db))
  r.GET("/user/:username", middleware.AuthRequest(controllers.GetUser(db)))

  log.Fatal(http.ListenAndServe(":" + constants.Port, r))
}
