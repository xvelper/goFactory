package database

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Username   string `gorm:"unique;not null"`
	Firstname  string
	Lastname   string
	Avatarpath string
	Email      string `gorm:"not null"`
	Password   string `gorm:"not null"`
	IsAdmin    bool   `gorm:"not null;default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Repository модель репозитория
type Repository struct {
	ID        uint   `gorm:"primaryKey"`
	UUID      string `gorm:"type:varchar(36);not null;unique"`
	Name      string `gorm:"not null"`
	OwnerID   uint   `gorm:"not null"`
	Owner     User   `gorm:"foreignKey:OwnerID"`
	Path      string `gorm:"not null"`
	IsPublic  bool   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Permission модель разрешения
type Permission struct {
	ID           uint       `gorm:"primaryKey"`
	UserID       uint       `gorm:"not null"`
	User         User       `gorm:"foreignKey:UserID"`
	RepositoryID uint       `gorm:"not null"`
	Repository   Repository `gorm:"foreignKey:RepositoryID"`
	Permission   string     `gorm:"not null"` // Например, "read", "write", "admin"
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	dsn := "host=localhost user=git password=198771 dbname=git port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&User{}, &Repository{}, &Permission{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

// HashPassword хеширование пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash проверка пароля
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
