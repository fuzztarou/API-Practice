package structure

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
