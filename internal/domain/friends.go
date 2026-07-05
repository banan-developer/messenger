package domain

type FriendResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type FriendRequest struct {
	FriendID int `json:"friend_id"`
}
