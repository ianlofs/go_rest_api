package controllers

import (
  "fmt"
  "log"
  "net/http"
  "time"

  "golang.org/x/crypto/bcrypt"
  "github.com/julienschmidt/httprouter"
  jwt "github.com/dgrijalva/jwt-go"

  "github.com/ianlofs/go_rest_api/constants"
  "github.com/ianlofs/go_rest_api/database"
  "github.com/ianlofs/go_rest_api/models"
)

func Login(db database.DB) httprouter.Handle {
  return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    username, password, ok := r.BasicAuth()
    if ok != true {
      http.Error(w, "Malformed credentials. HTTP Basic-Auth required.", http.StatusBadRequest)
      return
    }

    // authorize user
    user, err := models.FindUserByUsername(username, db)
    if err != nil {
      log.Println(err)
      http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
      return
    }
    err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
    if err != nil {
      log.Println(err)
      http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
      return
    }

    // create auth token
    token := jwt.New(jwt.GetSigningMethod("RS256"))
    token.Claims["ID"] = user.ID.Int64
    token.Claims["exp"] = time.Now().Add(72 * time.Hour).Unix() // Make the token valid for 3 days
    tokenString, err := token.SignedString(constants.PrivateKey)
    if err != nil {
      log.Println(err)
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }

    // craft response
    respBody, err := UserJSON{}.MarshalUser(user)
    if err != nil {
      log.Println(err)
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }

    // send response
    w.Header().Add("Bearer-Token", tokenString)
    w.Header().Set("Content-Type", "application/json")
    w.Write(respBody)
  }
}

func AuthRequest(handleFunc httprouter.Handle) httprouter.Handle {
  return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface {}, error) {
      if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
      }
      return constants.PublicKey, nil
    })

    if token == nil {
      http.Error(w, "Error: " + err.Error(), http.StatusBadRequest)
      return
    }

    if token.Valid {
      handleFunc(w,r, params);
    } else if ve, ok := err.(*jwt.ValidationError); ok {
      if ve.Errors & jwt.ValidationErrorMalformed != 0 {
        http.Error(w, "Error: Malformed auth token.", http.StatusBadRequest)
      } else if ve.Errors & (jwt.ValidationErrorExpired | jwt.ValidationErrorNotValidYet) != 0 {
        w.Header().Add("WWW-Authenticate", "Bearer")
        http.Error(w, "Error: Auth token expired.", http.StatusUnauthorized)
      } else {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      }
      return
    } else {
      http.Error(w, "Error: " + http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }
  }
}
