// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package rpchttpjson

import (
	"io"
	"net/rpc/jsonrpc"
)

// Структура запроса JSON RPC
type rpcRequest struct {
	r  io.Reader
	w  io.Writer
}

// Чтение данных
func (self *rpcRequest) Read(data []byte) (n int, err error) {
	return self.r.Read(data)
}

// Запись данных во
func (self *rpcRequest) Write(data []byte) (n int, err error) {
	return self.w.Write(data)
}

// Из-за поддержки интерфейса ReadWriteCloser введен этот пустой метод
func (self *rpcRequest) Close() error {
	return nil
}

// Исполнение JSON RPC запроса
func (self *rpcRequest) Serve() {
	jsonrpc.ServeConn(self)
}

// Создание экземпляра запроса
func newRequest(r io.Reader, w io.Writer) *rpcRequest {
	return &rpcRequest{ r, w }
}