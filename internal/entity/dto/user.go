package dto

type UserUpdateRequestDto struct {
	Login string `json:"login" binding:"omitempty,valid_login"`
}

type UserAdminUpdateRequestDto struct {
	IsAdmin   *bool `json:"is_admin,omitempty"`
	IsBlocked *bool `json:"is_blocked,omitempty"`
}
