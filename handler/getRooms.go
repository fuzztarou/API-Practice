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

func GetRooms(c echo.Context) error {
	//スライスを生成
	var rooms structure.RoomSlice
	//構造体の初期化
	r := structure.Room{
		Id:   0,
		Name: "",
	}
	//クエリ
	stmt, err := db.Db.Prepare("SELECT room_id, room_name FROM room")
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
