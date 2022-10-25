package handler

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/db"
	"api_test/model"
)

func CreateRoom(c echo.Context) error {
	//構造体の初期化
	roomStruct := new(model.Room)

	// リクエストのJSONから値を取得
	err := c.Bind(roomStruct)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.FailedToBindRequest,
		)
	}

	// リクエストJSONのnameが空だった時
	if roomStruct.Name == "" {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrRoomNameEmpty,
		)
	}

	// ルーム登録クエリ
	stmt_insert, err := db.Db.Prepare(
		"INSERT INTO room(room_name) VALUES(?)",
	)

	// クエリの実行
	res, err := stmt_insert.Exec(roomStruct.Name)
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

	// 登録確認クエリ
	stmt_confirm, err := db.Db.Prepare(
		"SELECT room_id, room_name FROM room WHERE room_id = ?",
	)

	// クエリの実行
	err = stmt_confirm.QueryRow(lastId).Scan(
		&roomStruct.Id,
		&roomStruct.Name,
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToGetLastInsertedRecord,
		)
	}

	return c.JSON(http.StatusOK, roomStruct)
}
