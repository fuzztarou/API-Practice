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

func CreateRoom(c echo.Context) error {
	//構造体の初期化
	r := new(structure.Room)
	// リクエストのJSONから値を取得
	err := c.Bind(r)
	if err != nil {
		log.Println("Error occured at Bind")
		log.Fatal(err)
	}
	// クエリ
	stmt, err := db.Db.Prepare("INSERT INTO room(room_name) VALUES(?)")
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
	stmt_confirm, err := db.Db.Prepare("SELECT room_id, room_name FROM room WHERE room_id = ?")
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
