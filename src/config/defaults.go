// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package config

const (
	// Значения по умолчанию настроек приложения
	DEF_APP_HOST     = "127.0.0.1"
	DEF_APP_PORT     = uint16(8080)
	DEF_APP_API_PATH = "/api"

	// Значения по умолчанию для подключения к MySQL БД
	DEF_DB_HOST      = "127.0.0.1"
	DEF_DB_PORT      = uint16(3306)
	DEF_DB_USER      = "root"
	DEF_DB_PASS      = ""
	DEF_DB_NAME      = "swetask"
	DEF_DB_COLT      = "utf8_general_ci"
)