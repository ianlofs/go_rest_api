package controllers

import (
  "encoding/json"
  "net/http"
  "log"

  "github.com/julienschmidt/httprouter"

  "github.com/ianlofs/go_rest_api/database"
  "github.com/ianlofs/go_rest_api/models"
)

type UserJSON struct {
  ID  int64 `json:"id"`
  Name string `json:"name"`
  Email string `json:"email"`
  Username string `json:"username"`
}

func (js UserJSON) MarshalUser(user *models.User) ([]byte, error){
  js.ID = user.ID.Int64
  js.Name = user.Name.String
  js.Email = user.Email.String
  js.Username = user.Username.String
  return json.Marshal(js)
}

func GetUser(db database.DB) (httprouter.Handle) {
  return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    id := params.ByName("id")
    user, err := models.FindUserByID(id, db)
    if err != nil {
      log.Println(err)
      http.Error(w, "Error: " + http.StatusText(http.StatusNotFound), http.StatusNotFound)
      return
    }

    respBody, err := UserJSON{}.MarshalUser(user)
    if err != nil {
      log.Println(err)
      http.Error(w, "Error: " + http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
      return
    }

    // send response
    w.Header().Set("Content-Type", "application/json")
    w.Write(respBody)
  }
}
