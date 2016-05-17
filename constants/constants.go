package constants

import (
  "fmt"
  "io/ioutil"
  "log"
  "os"
)


var (
  Envitonment string
  DB_URL string
  Port string
  PrivateKey []byte
  PublicKey []byte
)



func InitConstants() {
  // initialize environment neutral vars
  Envitonment = os.Getenv("ENVIRONMENT")
  Port = os.Getenv("PORT")

  DBUser := os.Getenv("DB_USER")
  DBPassword := os.Getenv("DB_PSSWRD")
  DBName := os.Getenv("DB_NAME")
  DBHost := os.Getenv("DB_HOST")
  DBPort := os.Getenv("DB_PORT")

  // set prod or dev specific vars
  var err error
  if Envitonment == "dev" {
    DB_URL = fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
      DBHost, DBUser, DBName, DBPassword, DBPort)
    PrivateKey, err = ioutil.ReadFile("constants/private.pem")
    PublicKey, err = ioutil.ReadFile("constants/pubkey.pem")
    if err != nil {
      log.Fatalln("Keys not read properly")
    }
  } else if Envitonment == "test" {
    // no tests yet
  } else {
    // env is prod
    DB_URL = os.Getenv("DATABASE_URL")
    PrivateKey = []byte(os.Getenv("PRIVATE_KEY"))
    PublicKey = []byte(os.Getenv("PUBLIC_KEY"))
  }
}
