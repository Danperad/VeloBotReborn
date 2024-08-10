package db_models

type Raiting struct {
	ResultId    int     `gorm:"column: result_id; primaryKey"`
	UserId      User    `gorm:"foreignKey: user_id"`
	MaxSpeed    float64 `gorm:"column: max_speed; type: numeric(10,5)"`
	Distance    float64 `gorm:"column: distance; type: numeric(10,5)"`
	SumDistance float64 `gorm:"column: sum_distance; type: numeric(10,5)"`
}
