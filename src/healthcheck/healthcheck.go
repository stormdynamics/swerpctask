// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package healthcheck

import (
	"log"
	"fmt"
	"net/http"
	"database/sql"
	"encoding/json"
)

// Структура обработчика очень простой проверки состояния сервиса 
type HealthCheck struct {
	db  *sql.DB
}

// Обработчик проверки состояния сервиса 
func (self *HealthCheck) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := self.db.Ping(); err != nil {
		fmt.Fprintf(w, `{"db":true, "db_error":"%s", "app":true}`,
			escapeJsonStr(err.Error()))
	} else {
		fmt.Fprint(w, `{"db":true, "app":true}`)
	}
}

// Функция преобразует спец. и управляющие символы к JSON формату,
// если преобразование закончилось ошибкой то вернем пустую строку
// и сделаем запись в лог
func escapeJsonStr(str string) string {
	if bytestr, err := json.Marshal(str); err != nil {
		log.Printf("error json string marshaling: %s\n", err.Error())
		return ""
	} else {
		str = string(bytestr)
		return str[1:len(str)-1]
	}
}

func NewHealthCheck(db *sql.DB) *HealthCheck {
	return &HealthCheck{ db }
}