package handler

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/db"
	"api_test/model"
)

func CreateUser(c echo.Context) error {
	//構造体初期化
	userStruct := new(model.User)

	// リクエストのJSONから値を取得
	err := c.Bind(userStruct)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.FailedToBindRequest,
		)
	}

	// リクエストJSONのnameが空だった時
	if userStruct.Name == "" {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrUserNameEmpty,
		)
	}

	// 新規登録クエリ
	stmt_create, err := db.Db.Prepare(
		"INSERT INTO user(user_name) VALUES(?)",
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToPrepareQuery,
		)
	}

	// クエリの実行
	res, err := stmt_create.Exec(userStruct.Name)
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
		"SELECT user_id, user_name FROM user WHERE user_id = ?",
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToPrepareQuery,
		)
	}

	// クエリの実行
	err = stmt_confirm.QueryRow(lastId).Scan(
		&userStruct.Id,
		&userStruct.Name,
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToGetLastInsertedRecord,
		)
	}

	return c.JSON(http.StatusOK, userStruct)
}
