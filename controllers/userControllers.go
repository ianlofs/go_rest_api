package controllers

import (
  "encoding/json"
  "io"
  "log"
  "net/http"
  "time"

  jwt "github.com/dgrijalva/jwt-go"
  "github.com/julienschmidt/httprouter"
  "golang.org/x/crypto/bcrypt"

  "github.com/ianlofs/go_rest_api/constants"
  "github.com/ianlofs/go_rest_api/database"
  "github.com/ianlofs/go_rest_api/models"
)

type GetUserJSON struct {
  ID  string `json:"id"`
  Name string `json:"name"`
  Email string `json:"email"`
  Username string `json:"username"`
}

type RegisterUserJSON struct {
  Name string `json:"name"`
  Email string `json:"email"`
  Username string `json:"username"`
  Password string  `json:"password"`
}


func marshalUser(user *models.User) ([]byte, error){
  var js GetUserJSON
  js.ID = user.ID.String
  js.Name = user.Name.String
  js.Email = user.Email.String
  js.Username = user.Username.String
  return json.Marshal(js)
}

func unMarshalUser(postBody io.ReadCloser) (*RegisterUserJSON, error){
  decoder := json.NewDecoder(postBody)
	var userJson RegisterUserJSON
	err := decoder.Decode(&userJson)
	if err != nil {
		return nil, err
	}
  return &userJson, nil
}

func Register(db database.DB) httprouter.Handle {
  return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    userJson, err := unMarshalUser(r.Body)
    if err != nil {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }

    // TODO: These two db calls should happen concurrently
    _, err = models.FindUserByUsername(userJson.Username, db)
    if err != nil {
      http.Error(w, "Error: Account with that username already exists", http.StatusOK)
      return
    }
    _, err = models.FindUserByEmail(userJson.Email, db)
    if err != nil {
      http.Error(w, "Error: Account with that email already exists", http.StatusOK)
      return
    }
    // hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userJson.Password), bcrypt.DefaultCost)
    if err != nil {
      http.Error(w, "Error: " + http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }
    // insert user
    err = models.InsertUser(db, userJson.Name, userJson.Email, userJson.Username, string(hashedPassword))
    if err != nil {
      http.Error(w, "Error: " + http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }
    w.WriteHeader(http.StatusCreated)
  }
}

func Login(db database.DB) httprouter.Handle {
  return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    username, password, ok := r.BasicAuth()
    if ok != true {
      http.Error(w, "Error: Malformed credentials. HTTP Basic-Auth required.", http.StatusBadRequest)
      return
    }
    // authorize user
    user, err := models.FindUserByUsername(username, db)
    if err != nil {
      log.Println(err)
      w.Header().Add("WWW-Authenticate", "Basic")
      http.Error(w, "Error: " + http.StatusText(http.StatusUnauthorized) + ". Username or password incorrect.", http.StatusUnauthorized)
      return
    }
    err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
    if err != nil {
      log.Println(err)
      w.Header().Add("WWW-Authenticate", "Basic")
      http.Error(w, "Error: " + http.StatusText(http.StatusUnauthorized) + ". Username or password incorrect.", http.StatusUnauthorized)
      return
    }
    // create auth token
    token := jwt.New(jwt.GetSigningMethod("RS256"))
    token.Claims["ID"] = user.ID.String
    token.Claims["exp"] = time.Now().Add(72 * time.Hour).Unix() // Make the token valid for 3 days
    tokenString, err := token.SignedString(constants.PrivateKey)
    if err != nil {
      log.Println(err)
      http.Error(w, "Error: " + http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }
    // craft response
    w.Header().Add("Bearer-Token", tokenString)
    returnUser(w, user);
  }
}

func GetUser(db database.DB) (httprouter.Handle) {
  return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    username := params.ByName("username")
    user, err := models.FindUserByUsername(username, db)
    if err != nil {
      log.Println(err)
      http.Error(w, "Error: " + http.StatusText(http.StatusNotFound), http.StatusNotFound)
      return
    }
    // craft response
    returnUser(w, user)
  }
}

func returnUser(w http.ResponseWriter, user *models.User) {
  respBody, err := marshalUser(user)
  if err != nil {
    log.Println(err)
    http.Error(w, "Error: " + http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }
  // send response
  w.Header().Set("Content-Type", "application/json")
  w.Write(respBody)
}
