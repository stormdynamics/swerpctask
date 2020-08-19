// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package config

import (
	"os"
	"log"
	"strconv"
)

// Основная структура конфигурации
type Config struct {
	App appConfig
	Db  dbConfig
}

// Структура для хранения значений конфигурации приложения
type appConfig struct {
	Host    string
	Port    uint16
	ApiPath string
}

// Структура для хранения значений конфигурации базы данных
type dbConfig struct {
	Host string
	Port uint16
	User string
	Pass string
	Name string
	Colt string
}

// Получение строкового значения переменной окружения с установкой значения
// по умочланию в случае возврата пустой строки
func getEnvString(name, defval string) string {
	env := os.Getenv(name)
	
	if env == "" {
		return defval
	}

	return env
}

// Получение целого беззнакового 2х байтовго значения переменной окружения
// с установкой значения по умочланию в случае возврата пустого значения
func getEnvUint16(name string, defval uint16) uint16 {
	env := os.Getenv(name)
	
	if env == "" {
		return defval
	}

	val, err := strconv.ParseUint(env, 10, 16)

	// если возникает ошибка в преобразовании строки переменной окружения
	// в целое беззнаковое 2х байтовое, то падаем с ошибкой в лог
	if err != nil {
		log.Fatalf("invalid environment variable value %s=%d, value must be a uint16", name, val)
	}

	if val == 0 {
		return defval
	}

	return uint16(val)
}

// Создание экземпляра конфигурации
func NewConfig() *Config {
		
	return &Config {

		App: appConfig {
			Host:    getEnvString("APP_HOST",     DEF_APP_HOST),
			Port:    getEnvUint16("APP_PORT",     DEF_APP_PORT),
			ApiPath: getEnvString("APP_API_PATH", DEF_APP_API_PATH),
		},

		Db: dbConfig {
			Host: getEnvString("DB_HOST", DEF_DB_HOST),
			Port: getEnvUint16("DB_PORT", DEF_DB_PORT),
			User: getEnvString("DB_USER", DEF_DB_USER),
			Pass: getEnvString("DB_PASS", DEF_DB_PASS),
			Name: getEnvString("DB_NAME", DEF_DB_NAME),
			Colt: getEnvString("DB_COLT", DEF_DB_COLT),
		},

	}
}