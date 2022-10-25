package handler

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/db"
	"api_test/model"
)

func PostChatToRoom(c echo.Context) error {
	//構造体の初期化
	chatStruct := new(model.Chat)

	// JSONから値を取得
	err := c.Bind(chatStruct)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.FailedToBindRequest,
		)
	}

	// 変数宣言
	user_id := chatStruct.UserId
	chat_txt := chatStruct.ChatTxt

	// room_idのパラメータ取得
	room_id := c.Param("id")

	// チャットルームにチャットを投稿するクエリ
	stmt_post, err := db.Db.Prepare(
		"INSERT INTO chat (user_id, room_id, chat_txt) VALUES (?,?,?)",
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToPrepareQuery,
		)
	}

	// 投稿しチャットを取得するクエリ
	stmt_confirm, err := db.Db.Prepare(
		"SELECT chat_id, user_name, room_name, chat_txt, created_at " +
			"FROM chat AS c " +
			"JOIN user AS u ON c.user_id = u.user_id " +
			"JOIN room AS r ON c.room_id = r.room_id " +
			"WHERE c.chat_id = ?",
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToPrepareQuery,
		)
	}

	// チャット投稿 クエリの実行
	res, err := stmt_post.Exec(user_id, room_id, chat_txt)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedInsertQuery,
		)
	}

	// InsertしたレコードのIDを取得
	lastId, err := res.LastInsertId()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToGetLastInsertedID,
		)
	}

	// チャット取得 クエリの実行
	err = stmt_confirm.QueryRow(lastId).Scan(
		&chatStruct.ChatId,
		&chatStruct.UserName,
		&chatStruct.RoomName,
		&chatStruct.ChatTxt,
		&chatStruct.CreatedAt,
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToGetLastInsertedRecord,
		)
	}

	return c.JSON(http.StatusOK, chatStruct)
}
