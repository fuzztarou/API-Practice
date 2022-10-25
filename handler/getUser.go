package handler

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/db"
	"api_test/model"
)

func GetUser(c echo.Context) error {
	//構造体初期化
	userStruct := new(model.User)
	// パラメータ取得
	id := c.Param("id")

	// クエリ
	stmt, err := db.Db.Prepare(
		"SELECT user_id, user_name FROM user WHERE user_id = ?",
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToPrepareQuery,
		)
	}

	// クエリの実行
	err = stmt.QueryRow(id).Scan(&userStruct.Id, &userStruct.Name)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToGetUserName,
		)
	}

	return c.JSON(http.StatusOK, userStruct)
}
