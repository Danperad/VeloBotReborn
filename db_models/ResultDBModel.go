package db_models

type Result struct {
	ResultId int     `gorm:"column: result_id; primaryKey"`
	UserId   User    `gorm:"foreignKey: user_id"`
	MaxSpeed float64 `gorm:"column: max_speed; type: numeric(10,5)"`
	Distance float64 `gorm:"column: distance; type: numeric(10,5)"`
}
