package dto

type TeamCreateRequestDto struct {
	Name        string  `json:"name" binding:"required,max=255"`
	Description *string `json:"description,omitempty"`
}

type TeamUpdateRequestDto struct {
	Name        string  `json:"name,omitempty" binding:"omitempty,max=255"`
	Description *string `json:"description,omitempty"`
}
