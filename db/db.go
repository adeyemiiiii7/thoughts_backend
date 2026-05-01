package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"thoughts_backend_api/models"
)

func InitDB(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	if err := db.AutoMigrate(
		&models.User{},
		&models.Follow{},
		&models.Interest{},
		&models.Thought{},
		&models.Comment{},
		&models.Reaction{},
		&models.EmailVerificationToken{},
		&models.PasswordResetToken{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
