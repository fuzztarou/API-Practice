package handler

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/db"
	"api_test/model"
)

func GetChat(c echo.Context) error {
	//スライスを作成
	var chats model.ChatSlice

	//構造体初期化
	chatStruct := model.Chat{
		ChatId:    "",
		UserId:    "",
		UserName:  "",
		RoomId:    "",
		RoomName:  "",
		ChatTxt:   "",
		CreatedAt: "",
		UpdatedAt: "",
	}
	// パラメータ取得
	id := c.Param("id")

	// クエリ
	stmt, err := db.Db.Prepare(
		"SELECT c.chat_id, u.user_name, c.chat_txt, c.updated_at " +
			"FROM chat AS c " +
			"JOIN user AS u ON c.user_id = u.user_id " +
			"WHERE c.room_id = ?",
	)
	if err != nil {
		log.Println(err)
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToPrepareQuery,
		)
	}

	//クエリの実行
	rows, err := stmt.Query(id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToGetRoomName,
		)
	}
	defer rows.Close()

	//クエリの実行結果をスライスに追加していく
	for rows.Next() {
		err = rows.Scan(&chatStruct.ChatId,
			&chatStruct.UserName,
			&chatStruct.ChatTxt,
			&chatStruct.UpdatedAt,
		)
		chats.Chat = append(chats.Chat, chatStruct)
	}
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToScanRoomIDandName,
		)
	}

	return c.JSON(http.StatusOK, chats)
}
