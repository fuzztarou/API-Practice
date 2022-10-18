package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/db"
	"api_test/structure"
)

func GetUser(c echo.Context) error {
	//構造体初期化
	u := new(structure.User)
	// パラメータ取得
	id := c.Param("id")
	// クエリ
	stmt, err := db.Db.Prepare("SELECT user_id, user_name FROM user WHERE user_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	// クエリの実行
	err = stmt.QueryRow(id).Scan(&u.Id, &u.Name)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, string(err.Error())+"\n")
	}
	//構造体をJSONに変換
	rtn_string, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, string(rtn_string)+"\n")
}
