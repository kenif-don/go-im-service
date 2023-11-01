package db

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/entity"
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
		panic(err)
	}
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		panic(err)
	}

	DefaultDB = &DB{
		Db: db,
	}
}

func NewDB() *DB {
	once.Do(func() {
		initDB()
	})
	return DefaultDB
}

type DB struct {
	Db *gorm.DB
}
