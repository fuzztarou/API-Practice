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

func CreateUser(c echo.Context) error {
	//構造体初期化
	u := new(structure.User)
	// リクエストのJSONから値を取得
	err := c.Bind(u)
	if err != nil {
		log.Println("Error occured at Bind")
		log.Fatal(err)
	}
	// 新規登録クエリ
	stmt_create, err := db.Db.Prepare("INSERT INTO user(user_name) VALUES(?)")
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
	stmt_confirm, err := db.Db.Prepare("SELECT user_id, user_name FROM user WHERE user_id = ?")
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
