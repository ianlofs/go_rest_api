package main

import (
  "log"
	"net/http"
  "os"

  "github.com/julienschmidt/httprouter"
  "github.com/joho/godotenv"

  "github.com/ianlofs/go_rest_api/controllers"
  "github.com/ianlofs/go_rest_api/database"
)

func main() {
  prod := os.Getenv("PRODUCTION")
  if prod != "true" {
    err := godotenv.Load()
    if err != nil {
      log.Fatal("Error loading .env file")
    }
  }

  db, err := database.New()
  if err != nil {
    log.Fatalln(err)
  }
  db.SetMaxOpenConns(40)
  defer db.Close()

  controllers.InitAuth()

  r := httprouter.New()
  r.POST("/login", controllers.Login(db))
  r.GET("/user/:id", controllers.AuthRequest(controllers.GetUser(db)))

  PORT := os.Getenv("PORT")

  if PORT == "" {
    PORT = "8080"
  }
  log.Fatal(http.ListenAndServe(":" + PORT, r))
}