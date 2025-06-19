package storage

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	client *gorm.DB
}

func NewMySQL(dsn string) (*MySQL, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	m := &MySQL{
		client: db,
	}

	err = m.Migrate()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *MySQL) Migrate() error {
	return m.client.AutoMigrate(&User{})
}

func (m *MySQL) CreateUser(user *User) error {
	return m.client.Create(user).Error
}

func (m *MySQL) FindUser(email string) (*User, error) {
	var user User

	err := m.client.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
