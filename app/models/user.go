package models

type User struct {
	Nickname string `json:"nickname,omitempty" gorm:"column:nickname;primaryKey"`
	Fullname string `json:"fullname,omitempty" validate:"required" gorm:"column:fullname"`
	About    string `json:"about,omitempty" gorm:"column:about"`
	Email    string `json:"email,omitempty" validate:"required" gorm:"column:email"`
}
