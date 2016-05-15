package authentication

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "time"

  "golang.org/x/crypto/bcrypt"
  jwt "github.com/dgrijalva/jwt-go"
  "github.com/julienschmidt/httprouter"

  "github.com/ianlofs/go_rest_api/database"
  "github.com/ianlofs/go_rest_api/models"
)


var (
  publicKey []byte
  privateKey []byte
)

func InitAuth() {
  initPrivateKey()
  initPublicKey()
}

func initPrivateKey() {
  var err error
	privateKey, err = ioutil.ReadFile("authentication/private.pem")

  // dont start if privatekey reading fails
  if err != nil {
    log.Fatalln("private key not read")
  }
}

func initPublicKey() {
  var err error
  publicKey, err = ioutil.ReadFile("authentication/pubkey.pem")

  // dont start if publickey reading fails
  if err != nil {
    log.Fatalln("public key not read")
  }
}


func AuthUser(db database.DB) (httprouter.Handle) {
  return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    password := r.Header.Get("Password")
    username := r.Header.Get("Username")

    user, err := models.FindUserByUsername(username, db)
    if err != nil {
      w.WriteHeader(http.StatusUnauthorized)
      w.Write([]byte(("ERROR: Bad credentials")))
      return
    }
    err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
    if err != nil {
      w.WriteHeader(http.StatusUnauthorized)
      w.Write([]byte(("ERROR: Bad credentials")))
      return
    }

    // Create a Token that will be signed with RSA 256.
    token := jwt.New(jwt.GetSigningMethod("RS256"))

    // The claims object allows you to store information in the actual token.
    token.Claims["ID"] = "123456789"
    token.Claims["exp"] = time.Now().Add(72 * time.Hour) // Make the token valid for 3 days
    var tokenString string
    tokenString, err = token.SignedString(privateKey)
    if err != nil {
      log.Println(err)
    }

    w.Header().Add("token", tokenString)
    w.WriteHeader(http.StatusOK)
  }
}

func AuthRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Println(r)
  token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface {}, error) {
    return publicKey, nil
  })
  if token.Valid {
    fmt.Println(token)
  } else {
    fmt.Println(err)
  }
}
