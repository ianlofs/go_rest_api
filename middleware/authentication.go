package middleware

import (
  "fmt"
  "net/http"

  "github.com/julienschmidt/httprouter"
  jwt "github.com/dgrijalva/jwt-go"

  "github.com/ianlofs/go_rest_api/constants"
)


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
