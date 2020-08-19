// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package rpchttpjson

import (
	"io"
	"log"
	"net/rpc"
	"net/http"
)

// Структура обработчика JSON RPC
type RpcHttpJson struct {
	api interface{}
}

// Обработчик JSON RPC
func (self *RpcHttpJson) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// Все ответы будут иметь формат JSON
	w.Header().Set("Content-Type", "application/json")

	// Создание нового запроса и его выполнение
	newRequest(req.Body, io.Writer(w)).Serve()
}

// Создание нового экземпляра структуры обработчика JSON RPC
// с регистрацией структуры и ее методов
func NewRpcHttpJson(api interface{}) *RpcHttpJson {

	err := rpc.Register(api)

	if err != nil {
		log.Fatal("error registering API", err)
	}

	return &RpcHttpJson{api}
}