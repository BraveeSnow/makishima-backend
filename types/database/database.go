package database

type User struct {
	ID           string `gorm:"primaryKey"`
	Username     string `gorm:"not null"`
	AccessToken  string `gorm:"unique;not null"`
	RefreshToken string `gorm:"unique;not null"`
	TokenExpiry  int64  `gorm:"not null"`
}
