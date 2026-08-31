package dto

type ProfileUpdateRequestDto struct {
	Name  string  `json:"name" binding:"required,max=255"`
	About *string `json:"about,omitempty"`
}
