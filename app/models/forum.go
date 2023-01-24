package models

type Forum struct {
	Title   string `json:"title,omitempty" gorm:"column:title"`
	User    string `json:"user,omitempty" gorm:"column:user_nickname"`
	Slug    string `json:"slug,omitempty" gorm:"column:slug"`
	Posts   int64  `json:"posts,omitempty" gorm:"column:posts"`
	Threads int64  `json:"threads,omitempty" gorm:"column:threads"`
}

type ForumUserRelation struct {
	Forum string `gorm:"column:forum"`
	User  string `gorm:"column:user_nickname"`
}

func (ForumUserRelation) TableName() string {
	return "forum_user"
}
