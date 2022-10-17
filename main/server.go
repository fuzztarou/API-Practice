package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/structure"
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
	e.POST("/room/:id/chat", postChatToRoom)

	// テスト
	e.POST("/test", testHandler)

	////サーバー起動////
	e.Logger.Fatal(e.Start(":1323"))
}

// 全チャットルームの取得
func getRooms(c echo.Context) error {
	//チャットルームの構造体
	// チャットルーム構造体を要素に持つスライス
	type roomslice struct {
		Room []structure.Room `json:"room"`
	}
	//スライスを生成
	var rooms roomslice
	//構造体の初期化
	r := structure.Room{
		Id:   0,
		Name: "",
	}
	//クエリ
	stmt, err := db.Prepare("SELECT room_id, room_name FROM room")
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
	//構造体初期化
	u := structure.User{
		Id:   0,
		Name: "",
	}
	// パラメータ取得
	id := c.Param("id")
	// クエリ
	stmt, err := db.Prepare("SELECT user_id, user_name FROM user WHERE user_id = ?")
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
	//構造体初期化
	u := new(structure.User)
	// リクエストのJSONから値を取得
	err = c.Bind(u)
	if err != nil {
		log.Println("Error occured at Bind")
		log.Fatal(err)
	}
	// 新規登録クエリ
	stmt_create, err := db.Prepare("INSERT INTO user(user_name) VALUES(?)")
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
	stmt_confirm, err := db.Prepare("SELECT user_id, user_name FROM user WHERE user_id = ?")
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
	//構造体の初期化
	r := new(structure.Room)
	// リクエストのJSONから値を取得
	err = c.Bind(r)
	if err != nil {
		log.Println("Error occured at Bind")
		log.Fatal(err)
	}
	// クエリ
	stmt, err := db.Prepare("INSERT INTO room(room_name) VALUES(?)")
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
	stmt_confirm, err := db.Prepare("SELECT room_id, room_name FROM room WHERE room_id = ?")
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
	// room_idをパラメータから取得
	room_id := c.Param("id")
	//構造体の初期化
	ur := new(structure.UserRoom)
	// リクエストのJSONから値を取得
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
		"SELECT count(*) FROM users_rooms WHERE user_id = ? AND room_id = ?",
	)
	// チャットルームにユーザーを登録するクエリ
	stmt_regist, err := db.Prepare(
		"INSERT INTO `users_rooms` (user_id, room_id) VALUES (?,?)",
	)
	// Insertしたレコードを取得するクエリ
	stmt_post_check, err := db.Prepare(
		"SELECT ur.user_room_id, u.user_name, r.room_name FROM users_rooms AS ur " +
			"JOIN user AS u ON ur.user_id = u.user_id " +
			"JOIN room AS r ON ur.room_id = r.room_id " +
			"WHERE ur.user_room_id = ?",
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
func postChatToRoom(c echo.Context) error {
	//構造体の初期化
	ch := new(structure.Chat)
	// room_idのパラメータ取得
	room_id := c.Param("id")
	// JSONから値を取得
	err := c.Bind(ch)
	// 変数宣言
	user_id := ch.UserId
	chat_txt := ch.ChatTxt
	rtn_string := ""

	// チャットルームにチャットを投稿するクエリ
	stmt_post, err := db.Prepare(
		"INSERT INTO chat (user_id, room_id, chat_txt) VALUES (?,?,?)",
	)
	// 投稿しチャットを取得するクエリ
	stmt_confirm, err := db.Prepare(
		"SELECT chat_id, user_name, room_name, chat_txt, created_at FROM chat AS c " +
			"JOIN user AS u ON c.user_id = u.user_id " +
			"JOIN room AS r ON c.room_id = r.room_id " +
			"WHERE c.chat_id = ?",
	)
	// チャット投稿 クエリの実行
	res, err := stmt_post.Exec(user_id, room_id, chat_txt)
	if err != nil {
		log.Println("Error occured at Exec")
		log.Fatal(err)
	}
	// InsertしたレコードのIDを取得
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	// チャット取得 クエリの実行
	err = stmt_confirm.QueryRow(lastId).Scan(&ch.ChatId, &ch.UserName, &ch.RoomName, &ch.ChatTxt, &ch.CreatedAt)
	if err != nil {
		log.Println("Error occured at SELECT")
		log.Fatal(err)
	}
	//構造体をJSONに変換
	json_string, err := json.Marshal(ch)
	if err != nil {
		log.Fatal(err)
	}
	rtn_string = string(json_string)

	return c.String(http.StatusOK, string(rtn_string)+"\n")
}

// テストハンドラ
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
