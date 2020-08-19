// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package config

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// Функция получения адреса вида <хост>:<порт> приложения
func (self *Config) AppHostPort() string {
	return fmt.Sprintf("%s:%d", self.App.Host, self.App.Port)
}

// Функция получения адреса вида <хост>:<порт> подключения к БД
func (self *Config) DbHostPort() string {
	return fmt.Sprintf("%s:%d", self.Db.Host, self.Db.Port)
}

// Функция формирования DSN строки для подключения к БД.
// Используется структура mysql.Config и функция FormatDSN
// из MySQL драйвера `github.com/go-sql-driver/mysql`
func (self *Config) DbDSN() string {
	// заполнение структуры
	dbcfg := mysql.Config {
		Net:                  "tcp",
		Addr:                 self.DbHostPort(),
		DBName:               self.Db.Name,
		User:                 self.Db.User,
		Passwd:               self.Db.Pass,
		Collation:            self.Db.Colt,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	// формирование строки DSN
	return dbcfg.FormatDSN()
}

// Функция получения значения API пути приложения
func (self *Config) AppApiPath() string {
	return self.App.ApiPath
}