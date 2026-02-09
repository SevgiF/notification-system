package mysql

import (
	"database/sql"
	env "github.com/SevgiF/notification-system/pkg/environment"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type MysqlConnectionManager struct {
	DB *sql.DB
}

// NewMySQLConnectionManager, ortam değişkenlerinden MySQL bağlantısını başlatır ve bir MySQLConnectionManager döner.
func NewMySQLConnectionManager() *MysqlConnectionManager {
	db := initMysqlConnection()
	return &MysqlConnectionManager{DB: db}
}

func initMysqlConnection() *sql.DB {

	mysqlConfig := mysql.Config{
		User:                 env.GetEnvOrFail("MYSQL_USER"),
		Passwd:               env.GetEnvOrFail("MYSQL_PASS"),
		Addr:                 env.GetEnvOrFail("MYSQL_ADDR"),
		DBName:               env.GetEnvOrFail("MYSQL_DB"),
		Net:                  "tcp",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		log.Fatalf("MySQL bağlantı hatası: %v", err)
	}

	// Bağlantı havuzu ayarları
	db.SetMaxOpenConns(env.GetIntEnv("MYSQL_MAX_OPEN_CONNS", 10))                         // Maksimum açık bağlantı sayısı
	db.SetMaxIdleConns(env.GetIntEnv("MYSQL_MAX_IDLE_CONNS", 5))                          // Maksimum boşta bekleyen bağlantı sayısı
	db.SetConnMaxLifetime(env.GetDurationEnv("MYSQL_CONN_MAX_LIFETIME", time.Hour))       // Bağlantı ömrü
	db.SetConnMaxIdleTime(env.GetDurationEnv("MYSQL_CONN_MAX_IDLE_TIME", time.Minute*15)) // Bağlantı boşta bekleme süresi

	// Bağlantının canlı olup olmadığını kontrol et
	if err := db.Ping(); err != nil {
		log.Fatalf("MySQL veritabanına erişilemiyor: %v", err)
	}

	return db

}
