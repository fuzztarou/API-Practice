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
	e.GET("/room", getRooms)
	//ユーザー取得
	e.GET("/user/:id", getUser)
	//ユーザー作成
	e.POST("/user", createUser)
	//ルーム作成
	e.POST("/room", createRoom)
	//ルームにユーザーを登録
	e.POST("/room/:id/user", registUserToRoom)
	//ルームにチャットを投稿
	//e.POST("/room/:id/chat", postChatToRoom)

	// テスト
	e.POST("/test", testHandler)

	////サーバー起動////
	e.Logger.Fatal(e.Start(":1323"))
}

// 全チャットルームの取得
func getRooms(c echo.Context) error {
	//チャットルームの構造体
	type room struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	// チャットルーム構造体を要素に持つスライス
	type roomslice struct {
		Room []room `json:"room"`
	}
	//スライスを生成
	var rooms roomslice
	//構造体の初期化
	r := room{
		Id:   0,
		Name: "",
	}
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
	rtn_string, err := json.Marshal(rooms)
	if err != nil {
		log.Fatal(err)
	}
	//jsonをリターン
	return c.String(http.StatusOK, string(rtn_string)+"\n")
}

// IDを指定してユーザー名の取得
func getUser(c echo.Context) error {
	//ユーザー構造体
	type user struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	//構造体初期化
	u := user{
		Id:   0,
		Name: "",
	}
	// パラメータ取得
	id := c.Param("id")
	// クエリ
	stmt, err := db.Prepare("select user_id, user_name from user where user_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	// クエリの実行
	err = stmt.QueryRow(id).Scan(&u.Id, &u.Name)
	if err != nil {
		log.Fatal(err)
	}
	//構造体をJSONに変換
	rtn_string, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, string(rtn_string)+"\n")
}

// ユーザーの新規登録
func createUser(c echo.Context) error {
	//ユーザー構造体
	type user struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	//構造体初期化
	u := new(user)
	// リクエストのJSONから値を取得
	err = c.Bind(u)
	if err != nil {
		log.Println("Error occured at Bind")
		log.Fatal(err)
	}
	// 新規登録クエリ
	stmt_create, err := db.Prepare("insert into user(user_name) values(?)")
	// クエリの実行
	res, err := stmt_create.Exec(u.Name)
	if err != nil {
		log.Println("Error occured at Exec")
		log.Fatal(err)
	}
	// InsertしたレコードのIDを取得
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	// 登録確認クエリ
	stmt_confirm, err := db.Prepare("select user_id, user_name from user where user_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	// クエリの実行
	err = stmt_confirm.QueryRow(lastId).Scan(&u.Id, &u.Name)
	if err != nil {
		log.Fatal(err)
	}
	//構造体をJSONに変換
	rtn_string, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, string(rtn_string)+"\n")
}

// チャットルームの新規登録
func createRoom(c echo.Context) error {
	//チャットルームの構造体
	type room struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	//構造体の初期化
	r := new(room)
	// リクエストのJSONから値を取得
	err = c.Bind(r)
	if err != nil {
		log.Println("Error occured at Bind")
		log.Fatal(err)
	}
	// クエリ
	stmt, err := db.Prepare("insert into room(room_name) values(?)")
	// クエリの実行
	res, err := stmt.Exec(r.Name)
	if err != nil {
		log.Println("Error occured at Exec")
		log.Fatal(err)
	}
	// InsertしたレコードのIDを取得
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	// 登録確認クエリ
	stmt_confirm, err := db.Prepare("select room_id, room_name from room where room_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	// クエリの実行
	err = stmt_confirm.QueryRow(lastId).Scan(&r.Id, &r.Name)
	if err != nil {
		log.Fatal(err)
	}
	//構造体をJSONに変換
	rtn_string, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, string(rtn_string)+"\n")
}

// ユーザーをチャットルームに登録
func registUserToRoom(c echo.Context) error {
	//ルームのユーザー登録情報の構造体
	type user_room struct {
		UserRoomId string `json:"id"`
		UserID     string `json:"user_id"`
		UserName   string `json:"user_name"`
		RoomID     string `json:"room_id"`
		RoomName   string `json:"room_name"`
	}
	// room_idをパラメータから取得
	room_id := c.Param("id")
	//構造体の初期化
	ur := new(user_room)
	/* 	ur := user_room{
	   		UserRoomId: "",
	   		UserID:     "",
	   		UserName:   "",
	   		RoomID:     room_id,
	   		RoomName:   "",
	   	}
	*/ // リクエストのJSONから値を取得
	err = c.Bind(ur)
	if err != nil {
		log.Println("Error occured at Bind")
		log.Fatal(err)
	}
	// user_idをJSONから取得
	user_id := ur.UserID
	// 構造体に値を代入
	ur.RoomID = room_id
	// 変数宣言
	count := ""
	rtn_string := ""

	// 登録があるかどうかをチェックするクエリ
	stmt_pre_check, err := db.Prepare(
		"select count(*) from users_rooms where user_id = ? and room_id = ?",
	)
	// チャットルームにユーザーを登録するクエリ
	stmt_regist, err := db.Prepare(
		"insert into `users_rooms` (user_id, room_id) values (?,?)",
	)
	// Insertしたレコードを取得するクエリ
	stmt_post_check, err := db.Prepare(
		"select ur.user_room_id, u.user_name, r.room_name from users_rooms as ur " +
			"join user as u on ur.user_id = u.user_id " +
			"join room as r on ur.room_id = r.room_id " +
			"where ur.user_room_id = ?",
	)

	// user_idとroom_idが一致するレコードをカウント
	err = stmt_pre_check.QueryRow(user_id, room_id).Scan(&count)
	if err != nil {
		log.Println("Error occured at Query")
		log.Fatal(err)
	}
	//登録済であれば新たに登録しない
	if count != "0" {
		rtn_string = "登録済です\n"
	} else {
		// 登録がなければ新規登録
		res, err := stmt_regist.Exec(user_id, room_id)
		if err != nil {
			log.Println("Error occured at Exec")
			log.Fatal(err)
		}
		//登録したレコードのIDを取得
		lastId, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		// 登録したレコードのuser_idとroom_idを取得
		err = stmt_post_check.QueryRow(lastId).Scan(&ur.UserRoomId, &ur.UserName, &ur.RoomName)
	}
	//構造体をJSONに変換
	json_string, err := json.Marshal(ur)
	if err != nil {
		log.Fatal(err)
	}
	rtn_string = string(json_string)

	return c.String(http.StatusOK, string(rtn_string)+"\n")
}

// ルームにチャットを投稿
/* func postChatToRoom(c echo.Context) error {
	//チャットの構造体
	type chat struct {
		ChatId    string `json:"chat_id"`
		RoomName  string `json:"room_name"`
		UserName  string `json:"user_name"`
		ChatText  string `json:"chat_text"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	//構造体の初期化
	ch := chat{
		ChatId:    "",
		RoomName:  "",
		UserName:  "",
		ChatText:  "",
		CreatedAt: "",
		UpdatedAt: "",
	}

	//curl -d 'user_id=2' -d 'text=test'  http://localhost:1323/room/1/chat
	// room_idのパラメータ取得
	room_id := c.Param("id")
	// user_idのフォーム値取得
	user_id := c.FormValue("user_id")
	chat_text := c.FormValue("text")
	// 変数宣言
	count := ""
	rtn_string := ""

	// チャットルームにチャットを投稿するクエリ
	stmt_post, err := db.Prepare(
		"insert into `users_rooms` (user_id, room_id) values (?,?)",
	)
	// Insertしたレコードを取得するクエリ
	stmt_post_check, err := db.Prepare(
		"select ur.user_room_id, u.user_name, r.room_name from users_rooms as ur " +
			"join user as u on ur.user_id = u.user_id " +
			"join room as r on ur.room_id = r.room_id " +
			"where ur.user_room_id = ?",
	)

	// user_idとroom_idが一致するレコードをカウント
	err = stmt_pre_check.QueryRow(user_id, room_id).Scan(&count)
	if err != nil {
		log.Println("Error occured at Query")
		log.Fatal(err)
	}
	//登録済であれば新たに登録しない
	if count != "0" {
		rtn_string = "登録済です\n"
	} else {
		// 登録がなければ新規登録
		res, err := stmt_regist.Exec(user_id, room_id)
		if err != nil {
			log.Println("Error occured at Exec")
			log.Fatal(err)
		}
		//登録したレコードのIDを取得
		lastId, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		// 登録したレコードのuser_idとroom_idを取得
		err = stmt_post_check.QueryRow(lastId).Scan(&ur.UserRoomId, &ur.UserName, &ur.RoomName)
	}
	//構造体をJSONに変換
	json_string, err := json.Marshal(ur)
	if err != nil {
		log.Fatal(err)
	}
	rtn_string = string(json_string)

	return c.String(http.StatusOK, string(rtn_string)+"\n")
}
*/
func testHandler(c echo.Context) error {
	type test struct {
		Id string `json:"id"`
	}

	t := new(test)
	if err = c.Bind(t); err != nil {
		log.Println("Error occured at Bind")
		log.Fatal(err)
	}
	test_id := t.Id
	log.Println(t)
	return c.String(http.StatusOK, test_id+"\n")
}
