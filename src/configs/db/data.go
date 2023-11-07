package db

import (
	"IM-Service/src/configs/conf"
	l "IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var (
	DefaultDB *DB
	once      sync.Once
	Ctx       context.Context
)

func initDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open(conf.DbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		l.Error(err)
	}
	err = db.AutoMigrate(&entity.User{}, &entity.FriendApply{}, &entity.Friend{}, &entity.Chat{})
	if err != nil {
		l.Error(err)
	}

	DefaultDB = &DB{
		Db: db,
	}
}

func NewDB() *DB {
	once.Do(func() {
		initDB()
		Ctx = context.Background()
	})
	return DefaultDB
}

type DB struct {
	Db *gorm.DB
}
