package db_models

type User struct {
	UserId   string `gorm:"column:user_id;type:varchar(150)"`
	UserName string `gorm:"column:user_name;type:varchar(511)"`
}
