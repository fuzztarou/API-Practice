package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/db"
	"api_test/handler"
	//"api_test/structure"
)

// エラー変数
var err error

// DBオブジェクトの生成
func init() {
	db.Db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/chat_app")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("DB Abstraction Succeeded!")
	}
}

func main() {
	//Echoオブジェクト？の生成
	e := echo.New()

	//
	//GET
	//
	//ルーム一覧を取得
	e.GET("/room", handler.GetRooms)
	//ユーザー取得
	e.GET("/user/:id", handler.GetUser)
	//ルームにチャットを投稿
	e.GET("/room/:id/chat", handler.GetChat)

	//
	//POST
	//
	//ユーザー作成
	e.POST("/user", handler.CreateUser)
	//ルーム作成
	e.POST("/room", handler.CreateRoom)
	//ルームにユーザーを登録
	e.POST("/room/:id/user", handler.RegistUserToRoom)
	//ルームにチャットを投稿
	e.POST("/room/:id/chat", handler.PostChatToRoom)

	////サーバー起動////
	e.Logger.Fatal(e.Start(":1323"))
}
