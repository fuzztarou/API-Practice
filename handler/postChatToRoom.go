package handler

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/db"
	"api_test/structure"
)

func PostChatToRoom(c echo.Context) error {
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
	stmt_post, err := db.Db.Prepare(
		"INSERT INTO chat (user_id, room_id, chat_txt) VALUES (?,?,?)",
	)
	// 投稿しチャットを取得するクエリ
	stmt_confirm, err := db.Db.Prepare(
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
