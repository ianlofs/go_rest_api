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
  return findUserByField("username", username, db)
}

func FindUserByID(id string, db database.DB) (*User, error) {
  return findUserByField("id", id, db)
}

func findUserByField(field string, value string, db database.DB) (*User, error) {
  user := User{}

  // definitely not safe from sql injection, will figure out a better way to do this
  sql := "SELECT id, name, email, username, password FROM users WHERE " + field + "=$1";
  stmt, err := db.Conn.Prepare(sql)
  if err != nil {
    return nil, err
  }
  defer stmt.Close()
  stmt.QueryRow(value).Scan(&user.ID, &user.Name, &user.Email, &user.Username, &user.Password)
  fmt.Println(user.Username.String)
  return &user, nil
}
