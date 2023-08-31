package models

type PhoneOutput struct {
	UserID      int    `json:"user_id"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
	IsFax       bool   `json:"is_fax"`
}

type Phone struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
	IsFax       bool   `json:"is_fax"`
}

type PhoneUpdateInput struct {
	Phone       string `json:"phone"`
	Description string `json:"description"`
	PhoneId     int    `json:"phone_id"`
	IsFax       bool   `json:"is_fax"`
}
