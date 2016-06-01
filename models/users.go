package models

import (
  "crypto/rand"
  "database/sql"
 	"fmt"
 	"time"

  "github.com/ianlofs/go_rest_api/database"
)

type User struct {
  ID       sql.NullString
  Name     sql.NullString
  Email    sql.NullString
  Username sql.NullString
  Password sql.NullString
  Activated sql.NullBool
}

func makeID () (string, error) {
  unix32bits := uint32(time.Now().UTC().Unix())
 	buff := make([]byte, 12)
 	numRead, err := rand.Read(buff)
 	if numRead != len(buff) || err != nil {
 		return "", err
 	}

 	return fmt.Sprintf("%x-%x-%x-%x-%x-%x", unix32bits, buff[0:2], buff[2:4], buff[4:6], buff[6:8], buff[8:]), nil
}

func InsertUser(db database.DB,  name string, email string, username string, hashedPassword string) error {
    sql := "INSERT INTO users(id, name, email, username, password, activated) VALUES($1, $2, $3, $4, $5, $6);"
    stmt, err := db.Conn.Prepare(sql)
    if err != nil {
      return err
    }
    defer stmt.Close()
    id, err := makeID()
    if err != nil {
      return err
    }
    _, err = stmt.Exec(id, name, email, username, hashedPassword, false)
    if err != nil {
      return err
    }
    return nil
}

func ActivateUser(db database.DB, id string) error {
  sql := "UDPATE users SET activated=$1 WHERE id=$2;"
  stmt, err := db.Conn.Prepare(sql)
  if err != nil {
    return err
  }
  defer stmt.Close()
  _, err = stmt.Exec(true, id)
  if err != nil {
    return err
  }
  return nil
}

func FindUserByEmail(email string, db database.DB) (*User, error) {
  selectData := map[string]interface{}{"email":email}
  return findUserByField( db, selectData)
}

func FindUserByEmailOrUsername(db database.DB, email, username) (*User, error) {
  selectData := map[string]interface{}{"email":email, "username": username}
  return findUserByField(db, selectData)
}

func FindUserByID(id string, db database.DB) (*User, error) {
  return findUserByField( db, ["id"], [id])
}

func findUserByField(db database.DB, selectData map[string]interface{}) (*User, error) {
  user := User{}
  sql := "SELECT id, name, email, username, password FROM users WHERE"
  for key, value := range(selectData) {
    if value != nil {
      sql += " " + key + "=" + value + " "
    }
  }
  sql += ";"
  stmt, err := db.Conn.Prepare(sql)
  if err != nil {
    return nil, err
  }
  defer stmt.Close()
  err = stmt.QueryRow(value).Scan(&user.ID, &user.Name, &user.Email, &user.Username, &user.Password)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  return &user, nil
}
