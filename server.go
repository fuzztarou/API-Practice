package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

// 変数
var (
	db  *sql.DB
	err error
)

// DBオブジェクトの生成
func init() {
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/chat_app")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("DB Abstraction Succeeded!")
	}
}

func main() {
	//Echoオブジェクト？の生成
	e := echo.New()
	////ルーティング////
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	//ユーザー取得
	e.GET("/user/:id", getUser)
	//ユーザー取作成
	e.POST("/create-user", createUser)
	////サーバー起動////
	e.Logger.Fatal(e.Start(":1323"))
}

func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	name := ""

	stmt, err := db.Prepare("select user_name from user where user_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	err = stmt.QueryRow(id).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(id, name)

	return c.String(http.StatusOK, name)
}

func createUser(c echo.Context) error {
	name := c.FormValue("name")

	stmt, err := db.Prepare("insert into user(user_name) values(?)")

	res, err := stmt.Exec(name)
	if err != nil {
		log.Println("Error occured at Exec")
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("New User %s is created! id=%d", name, lastId)

	return c.String(http.StatusOK, name)
}
