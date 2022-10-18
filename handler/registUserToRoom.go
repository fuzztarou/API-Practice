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

func RegistUserToRoom(c echo.Context) error {
	// room_idをパラメータから取得
	room_id := c.Param("id")
	//構造体の初期化
	ur := new(structure.UserRoom)
	// リクエストのJSONから値を取得
	err := c.Bind(ur)
	if err != nil {
		log.Println("Error occured at Bind")
		log.Fatal(err)
	}
	// user_idをJSONから取得
	user_id := ur.UserID
	// 構造体に値を代入
	ur.RoomID = room_id
	// 変数宣言
	count := ""
	rtn_string := ""

	// 登録があるかどうかをチェックするクエリ
	stmt_pre_check, err := db.Db.Prepare(
		"SELECT count(*) FROM users_rooms WHERE user_id = ? AND room_id = ?",
	)
	// チャットルームにユーザーを登録するクエリ
	stmt_regist, err := db.Db.Prepare(
		"INSERT INTO `users_rooms` (user_id, room_id) VALUES (?,?)",
	)
	// Insertしたレコードを取得するクエリ
	stmt_post_check, err := db.Db.Prepare(
		"SELECT ur.user_room_id, u.user_name, r.room_name FROM users_rooms AS ur " +
			"JOIN user AS u ON ur.user_id = u.user_id " +
			"JOIN room AS r ON ur.room_id = r.room_id " +
			"WHERE ur.user_room_id = ?",
	)

	// user_idとroom_idが一致するレコードをカウント
	err = stmt_pre_check.QueryRow(user_id, room_id).Scan(&count)
	if err != nil {
		log.Println("Error occured at Query")
		log.Fatal(err)
	}
	//登録済であれば新たに登録しない
	if count != "0" {
		rtn_string = "登録済です\n"
	} else {
		// 登録がなければ新規登録
		res, err := stmt_regist.Exec(user_id, room_id)
		if err != nil {
			log.Println("Error occured at Exec")
			log.Fatal(err)
		}
		//登録したレコードのIDを取得
		lastId, err := res.LastInsertId()
		if err != nil {
			log.Println("Last inserted ID: ", lastId)
			log.Fatal(err)
		}

		// 登録したレコードのuser_idとroom_idを取得
		err = stmt_post_check.QueryRow(lastId).Scan(&ur.UserRoomId, &ur.UserName, &ur.RoomName)
		if err != nil {
			log.Println("Error occured at Query Row")
			log.Fatal(err)
		}
		//構造体をJSONに変換
		json_string, err := json.Marshal(ur)
		if err != nil {
			log.Fatal(err)
		}
		rtn_string = string(json_string)
	}

	return c.String(http.StatusOK, string(rtn_string)+"\n")
}
