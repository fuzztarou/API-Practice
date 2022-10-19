package handler

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/db"
	"api_test/structure"
)

func CreateRoom(c echo.Context) error {
	//構造体の初期化
	roomStruct := new(structure.Room)
	//構造体の初期化
	errStruct := new(structure.Error)

	// リクエストのJSONから値を取得
	err := c.Bind(roomStruct)
	if err != nil {
		errStruct.Message = "Failed to bind your request"
		errStruct.ErrorCode = "100"
		return c.JSON(http.StatusBadRequest, errStruct)
	}
	// リクエストJSONのnameが空だった時
	if roomStruct.Name == "" {
		errStruct.Message = "Room name should not be empty"
		errStruct.ErrorCode = "100"
		return c.JSON(http.StatusBadRequest, errStruct)
	}

	// ルーム登録クエリ
	stmt_insert, err := db.Db.Prepare("INSERT INTO room(room_name) VALUES(?)")

	// クエリの実行
	res, err := stmt_insert.Exec(roomStruct.Name)
	if err != nil {
		errStruct.Message = "Query Failed at INSERT"
		errStruct.ErrorCode = "100"
		return c.JSON(http.StatusInternalServerError, errStruct)
	}

	// InsertしたレコードのIDを取得
	lastId, err := res.LastInsertId()
	if err != nil {
		errStruct.Message = "Failed to get last inserted ID"
		errStruct.ErrorCode = "100"
		return c.JSON(http.StatusInternalServerError, errStruct)
	}

	// 登録確認クエリ
	stmt_confirm, err := db.Db.Prepare("SELECT room_id, room_name FROM room WHERE room_id = ?")

	// クエリの実行
	err = stmt_confirm.QueryRow(lastId).Scan(&roomStruct.Id, &roomStruct.Name)
	if err != nil {
		errStruct.Message = "Failed to get last inserted record"
		errStruct.ErrorCode = "100"
		return c.JSON(http.StatusInternalServerError, errStruct)
	}

	return c.JSON(http.StatusOK, roomStruct)
}
