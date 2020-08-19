// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package rpcapi

import (
	"github.com/stormdynamics/swerpctask/src/incrementor"
)

// Структура RPC API
type Api struct {
	incr  *incrementor.Incrementor
}

// Структура для получения настроек инкрементатора
// в методе SetSettings
type ApiSettings struct {
	Maxval int `json:"max_val"`
	Incval int `json:"inc_val"`
}

// RPC API метод получение текущего значения счетчика
func (self *Api) GetNumber(_ interface{}, reply *int) error {
	*reply = self.incr.GetNumber()
	return nil
}

// RPC API метод увеличения счетчика
func (self *Api) IncrementNumber(_ interface{}, reply *int) error {
	return self.incr.IncrementNumber()
}

// RPC API метод установки максимального значения
func (self *Api) SetMaxValue(maxval *int, reply *int) error {
	return self.incr.SetMaximumValue(*maxval)
}

// RPC API метод установки значения приращения
func (self *Api) SetIncVal(incval *int, reply *int) error {
	return self.incr.SetIncrValue(*incval)
}

// RPC API метод установки настроек приложения
func (self *Api) SetSettings(settings *ApiSettings, reply *int) error {
	return self.incr.SetSettings(settings.Maxval, settings.Incval)
}

// Создание экземпляра структуры RPC API
func NewApi(incr *incrementor.Incrementor) *Api {
	return &Api{incr}
}