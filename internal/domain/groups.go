package domain

type CreateGroupRequest struct {
	Title   string `json:"title"`
	UserIDs []int  `json:"user_ids"`
}
