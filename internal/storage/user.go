package storage

type User struct {
	ID       string `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Password string
}
