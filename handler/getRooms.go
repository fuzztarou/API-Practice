package handler

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/db"
	"api_test/model"
)

func GetRooms(c echo.Context) error {
	//スライスを生成
	var rooms model.RoomSlice

	//構造体の初期化
	roomStruct := model.Room{
		Id:   0,
		Name: "",
	}

	//クエリ
	stmt, err := db.Db.Prepare("SELECT room_id, room_name FROM room")
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToPrepareQuery,
		)
	}

	//クエリの実行
	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToGetRoomName,
		)
	}
	defer rows.Close()

	//クエリの実行結果をスライスに追加していく
	for rows.Next() {
		err = rows.Scan(&roomStruct.Id, &roomStruct.Name)
		rooms.Room = append(rooms.Room, roomStruct)
	}
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToScanRoomIDandName,
		)
	}

	//jsonをリターン
	return c.JSON(http.StatusOK, rooms)
}
