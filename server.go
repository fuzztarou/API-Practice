package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"strconv"

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
	//ルーム一覧を取得
	e.GET("/", getRooms)
	//ユーザー取得
	e.GET("/user/:id", getUser)
	//ユーザー作成
	e.POST("/create-user", createUser)
	//ルーム作成
	e.POST("/create-room", createRoom)
	//ルームにユーザーを登録
	e.POST("/regist-user", registUserRoom)

	////サーバー起動////
	e.Logger.Fatal(e.Start(":1323"))
}

// 全チャットルームの取得
func getRooms(c echo.Context) error {
	//チャットルームの構造体
	type room struct {
		Id   int
		Name string
	}
	// チャットルーム構造体を要素に持つスライス
	type roomslice struct {
		Room []room
	}
	//構造体初期化
	r := room{}
	r.Id = 0
	r.Name = ""
	//スライスを生成
	var rooms roomslice
	//クエリ
	stmt, err := db.Prepare("select room_id, room_name from room")
	if err != nil {
		log.Fatal(err)
	}
	//クエリの実行
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	//クエリの実行結果をスライスに追加していく
	for rows.Next() {
		err = rows.Scan(&r.Id, &r.Name)
		rooms.Room = append(rooms.Room, r)
	}
	if err != nil {
		log.Fatal(err)
	}
	//スライスをjsonに変換
	output, err := json.Marshal(rooms)
	if err != nil {
		log.Fatal(err)
	}
	//jsonをリターン
	return c.String(http.StatusOK, string(output))
}

// IDを指定してユーザー名の取得
func getUser(c echo.Context) error {
	// パラメータ取得
	id := c.Param("id")
	// 変数宣言
	name := ""
	// クエリ
	stmt, err := db.Prepare("select user_name from user where user_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	// クエリの実行
	err = stmt.QueryRow(id).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(id, name)

	return c.String(http.StatusOK, name+"\n")
}

// ユーザーの新規登録
func createUser(c echo.Context) error {
	// パラメータ取得
	name := c.FormValue("name")
	// クエリ
	stmt, err := db.Prepare("insert into user(user_name) values(?)")
	// クエリの実行
	res, err := stmt.Exec(name)
	if err != nil {
		log.Println("Error occured at Exec")
		log.Fatal(err)
	}
	// InsertしたレコードのIDを取得
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("New User %s is created! id=%d", name, lastId)

	return c.String(http.StatusOK, "")
}

// チャットルームの新規登録
func createRoom(c echo.Context) error {
	// パラメータ取得
	name := c.FormValue("name")
	// クエリ
	stmt, err := db.Prepare("insert into room(room_name) values(?)")
	// クエリの実行
	res, err := stmt.Exec(name)
	if err != nil {
		log.Println("Error occured at Exec")
		log.Fatal(err)
	}
	// InsertしたレコードのIDを取得
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("New Room %s is created! id=%d", name, lastId)

	return c.String(http.StatusOK, "")
}

// ユーザーをチャットルームに登録
func registUserRoom(c echo.Context) error {
	// パラメータ取得
	user := c.FormValue("user")
	room := c.FormValue("room")
	// 変数宣言
	id := ""
	count := ""
	rtn_string := ""
	// 登録があるかどうかをチェックするクエリ
	stmt_pre_check, err := db.Prepare("select count(*) from users_rooms where user_id = ? and room_id = ?")
	// ユーザーとチャットルームを結びつけるクエリ
	stmt_regist, err := db.Prepare("insert into `users_rooms` (user_id, room_id) values (?,?)")
	// Insertしたレコードを取得するクエリ
	stmt_post_check, err := db.Prepare("select user_room_id, from users_rooms where user_id = ? and room_id = ?")

	// クエリの実行
	err = stmt_pre_check.QueryRow(user, room).Scan(&count)
	if err != nil {
		log.Println("Error occured at Query")
		log.Fatal(err)
	}
	//登録済であれば新たに登録しない
	if count != "0" {
		rtn_string = "登録済です\n"
	} else {
		// 登録がなければ新規登録
		_, err = stmt_regist.Exec(user, room)
		if err != nil {
			log.Println("Error occured at Exec")
			log.Fatal(err)
		}
		// 登録したルームのIDを取得
		err = stmt_post_check.QueryRow(user, room).Scan(&id)
		rtn_string = "登録完了  id=" + id + "\n"
	}
	return c.String(http.StatusOK, rtn_string)
}
