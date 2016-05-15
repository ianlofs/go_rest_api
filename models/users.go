package models

import (
  "database/sql"
  "fmt"

  "github.com/ianlofs/go_rest_api/database"
)

type User struct {
  ID       sql.NullInt64
  Name     sql.NullString
  Email    sql.NullString
  Username sql.NullString
  Password sql.NullString
}

func FindUserByUsername(username string, db database.DB) (*User, error) {

  fmt.Println(username)
  user := User{}
  stmt, err := db.Conn.Prepare("SELECT id, name, email, username, password FROM users WHERE username=$1;")
  if err != nil {
    return nil, err
  }
  defer stmt.Close()
  stmt.QueryRow(username).Scan(&user.ID, &user.Name, &user.Email, &user.Username, &user.Password)
  return &user, nil
}

func (user User) FindUserByEmail(email string, db database.DB) (*User, error) {

  return &User{}, nil
}
