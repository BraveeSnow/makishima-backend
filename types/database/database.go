package database

type User struct {
	ID           string `gorm:"primaryKey"`
	Username     string `gorm:"not null"`
	AccessToken  string `gorm:"unique;not null"`
	RefreshToken string `gorm:"unique;not null"`
	TokenExpiry  int64  `gorm:"not null"`
}

type AnilistUser struct {
	ID           int    `gorm:"primaryKey"`
	AccessToken  string `gorm:"unique;not null"`
	RefreshToken string `gorm:"unique;not null"`
	TokenExpiry  int64  `gorm:"not null"`

	// user foreign key
	UserID string `gorm:"not null"`
	User   User
}
