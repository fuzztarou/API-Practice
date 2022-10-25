package handler

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"api_test/db"
	"api_test/model"
)

func RegistUserToRoom(c echo.Context) error {
	// room_idをパラメータから取得
	room_id := c.Param("id")

	//構造体の初期化
	userRoomStruct := new(model.UserRoom)

	// リクエストのJSONから値を取得
	err := c.Bind(userRoomStruct)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.FailedToBindRequest)
	}

	// user_idをJSONから取得
	user_id := userRoomStruct.UserID
	// 構造体に値を代入
	userRoomStruct.RoomID = room_id
	// 変数宣言
	count := ""

	// 登録があるかどうかをチェックするクエリ
	stmt_pre_check, err := db.Db.Prepare(
		"SELECT count(*) FROM users_rooms WHERE user_id = ? AND room_id = ?",
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToPrepareQuery,
		)
	}

	// チャットルームにユーザーを登録するクエリ
	stmt_regist, err := db.Db.Prepare(
		"INSERT INTO `users_rooms` (user_id, room_id) VALUES (?,?)",
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToPrepareQuery,
		)
	}

	// Insertしたレコードを取得するクエリ
	stmt_post_check, err := db.Db.Prepare(
		"SELECT userRoomStruct.user_room_id, u.user_name, r.room_name " +
			"FROM users_rooms AS userRoomStruct " +
			"JOIN user AS u ON userRoomStruct.user_id = u.user_id " +
			"JOIN room AS r ON userRoomStruct.room_id = r.room_id " +
			"WHERE userRoomStruct.user_room_id = ?",
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToPrepareQuery,
		)
	}

	// user_idとroom_idが一致するレコードをカウント
	err = stmt_pre_check.QueryRow(user_id, room_id).Scan(&count)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.FailedToGetRequiredData,
		)
	}

	//登録済であれば新たに登録しない
	if count != "0" {
		return c.JSON(
			http.StatusInternalServerError,
			model.DataAlreadyExists,
		)
	} else {
		// 登録がなければ新規登録
		res, err := stmt_regist.Exec(user_id, room_id)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				model.FailedInsertQuery,
			)
		}
		//登録したレコードのIDを取得
		lastId, err := res.LastInsertId()
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				model.FailedToGetLastInsertedID,
			)
		}

		// 登録したレコードのuser_idとroom_idを取得
		err = stmt_post_check.QueryRow(lastId).Scan(
			&userRoomStruct.UserRoomId,
			&userRoomStruct.UserName,
			&userRoomStruct.RoomName,
		)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				model.FailedToGetLastInsertedRecord,
			)
		}
	}

	return c.JSON(http.StatusOK, userRoomStruct)
}
