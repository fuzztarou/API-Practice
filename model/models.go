package model

// echoのBinding  https://echo.labstack.com/guide/binding/

// ユーザー構造体
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// チャットルームの構造体
type Room struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// チャットルーム構造体を要素に持つスライス
type RoomSlice struct {
	Room []Room `json:"room"`
}

// ルームのユーザー登録情報の構造体
type UserRoom struct {
	UserRoomId string `json:"id"`
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name"`
	RoomID     string `json:"room_id"`
	RoomName   string `json:"room_name"`
}

// チャットの構造体
type Chat struct {
	ChatId    string `json:"chat_id"`
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	RoomId    string `json:"room_id"`
	RoomName  string `json:"room_name"`
	ChatTxt   string `json:"chat_txt"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// チャットルーム構造体を要素に持つスライス
type ChatSlice struct {
	Chat []Chat `json:"chat"`
}

// エラー構造体
type Error struct {
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}

// Bad Request系エラー
var FailedToBindRequest = Error{
	Message:   "Failed to bind your request",
	ErrorCode: "101",
}

var ErrRoomNameEmpty = Error{
	Message:   "Room name should not be empty",
	ErrorCode: "102",
}

var ErrUserNameEmpty = Error{
	Message:   "User name should not be empty",
	ErrorCode: "103",
}

// Query関係エラー
var FailedToPrepareQuery = Error{
	Message:   "Failed to prepare query",
	ErrorCode: "200",
}

// DB関係エラー
var FailedInsertQuery = Error{
	Message:   "Query Failed at INSERT",
	ErrorCode: "301",
}

var FailedToGetLastInsertedID = Error{
	Message:   "Failed to get last inserted ID",
	ErrorCode: "302",
}

var FailedToGetLastInsertedRecord = Error{
	Message:   "Failed to get last inserted record",
	ErrorCode: "303",
}

var FailedToGetRoomName = Error{
	Message:   "Failed to get room name(s) from DB",
	ErrorCode: "304",
}

var FailedToScanRoomIDandName = Error{
	Message:   "Failed to scan room ID and room name",
	ErrorCode: "305",
}

var FailedToGetUserName = Error{
	Message:   "Failed to get user name",
	ErrorCode: "306",
}

var FailedToGetRequiredData = Error{
	Message:   "Failed to get required data from DB",
	ErrorCode: "307",
}

var DataAlreadyExists = Error{
	Message:   "Data is already existed",
	ErrorCode: "308",
}
