// +build integration

package main

import (
	"testing"
	"os"
    "fmt"
	"bytes"
	"net/http"
    "database/sql"
	"net/http/httptest"
	"encoding/json"

	"github.com/stormdynamics/swerpctask/src/config"
	"github.com/stormdynamics/swerpctask/src/healthcheck"
	"github.com/stormdynamics/swerpctask/src/incrementor"
	"github.com/stormdynamics/swerpctask/src/rpchttpjson"
	"github.com/stormdynamics/swerpctask/src/rpcapi.v1"
)

// Структура для разбора JSON RPC ответа
type jsonRpcResult struct {
	Jsonrpc string      `json:"jsonrpc,omitempty"`
	Id      uint        `json:"id"`
	Error   interface{} `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

// Структура для формирования JSON RPC запроса
type jsonRpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  interface{}   `json:"params"`
	Id      uint          `json:"id"` 
}

var (
	cfg *config.Config
	db  *sql.DB
	inc *incrementor.Incrementor
	api *rpcapi.Api
	rpc *rpchttpjson.RpcHttpJson
)

// Выполняем первоначальные инициализации для тестов
func TestMain(m *testing.M) {
	var err error

	cfg      = config.NewConfig()
	db, err  = sql.Open("mysql", cfg.DbDSN())

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	inc = incrementor.NewIncrementor(db)
	api = rpcapi.NewApi(inc)
	rpc = rpchttpjson.NewRpcHttpJson(api)

	// т.к. по заданию у нас нет возможности выставить текущее значение
    // счетчика через API, то используем такой workaround
	err = inc.SetSettings(1, 1)
	
	if err != nil {
		panic(err)
	}

	err = inc.SetSettings(incrementor.MAX_INT, 1)

	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Exit(code)
}

// Тестирование обработчика проверки состояния сервиса 
func TestHealthCheck(t *testing.T) {
	hc := healthcheck.NewHealthCheck(db)

	req, err := http.NewRequest("GET", "/health-check", nil)

	if err != nil {
        t.Fatal(err)
    }

    recorder := httptest.NewRecorder()

    hc.ServeHTTP(recorder, req)

    if status := recorder.Code; status != http.StatusOK {
    	t.Errorf("/health-check wrong status code: %d want %d",
            status, http.StatusOK)
        return
    }

    hckres := struct {
		Db  bool `json:"db"` 
		App bool `json:"app"`
	}{ false, false }

    if err := json.Unmarshal([]byte(recorder.Body.String()), &hckres); err != nil {
    	t.Errorf("/health-check unexpected json body: %s", err.Error())
    	return
    }

    if !hckres.Db || !hckres.App {
    	t.Errorf("/health-check failed want db:true, app:true: db:%t app:%t",
            hckres.Db, hckres.App)
    }
}

// Тестирование RPC функции Api.GetNumber
func TestApiGetNumber(t *testing.T) {
    currentval, err := ApiGetNumber(123)

    if err != nil {
        t.Errorf(err.Error())
        return
    }

    if currentval != 0 {
		t.Errorf("rpc api invalid value: %v want 0", currentval)
    }
}

// Тестирование RPC функции Api.GetNumber
func TestApiIncrementNumber(t *testing.T) {

	// ===========================================
	// получение текущего значения счетчика
	// ===========================================

	currentval, err := ApiGetNumber(551)

    if err != nil {
        t.Errorf(err.Error())
        return
    }

    // ===========================================
    // инкрементация счетчика на 1
    // ===========================================

    if err := ApiIncrementNumber(552); err != nil {
        t.Errorf(err.Error())
        return
    }

    // ===========================================
    // получение новго значения счетчика
    // ===========================================

    currentnew, err := ApiGetNumber(551)

    if err != nil {
        t.Errorf(err.Error())
        return
    }

    // новое значение должно быть больше предыдущего на 1
    if currentval + 1 != currentnew {
    	t.Errorf("invalid new increment value %d want %d", currentnew, currentval + 1)
    	return
    }
}

// Тестирование RPC функции Api.SetSettings
func TestApiSetSettings(t *testing.T) {

	// ===========================================
    // установка максимального значения и значения приращения
    // ===========================================

    if err := ApiSetSettings(155, 33, 3); err != nil {
        t.Error(err.Error())
        return
    }

    // ===========================================
	// получение текущего значения счетчика
	// ===========================================

	currentval, err := ApiGetNumber(641)

    if err != nil {
        t.Errorf(err.Error())
        return
    }

    // ===========================================
    // инкрементация счетчика на 3
    // ===========================================

    if err := ApiIncrementNumber(771); err != nil {
        t.Errorf(err.Error())
        return
    }

    // ===========================================
    // получение новго значения счетчика
    // ===========================================

    currentnew, err := ApiGetNumber(512)

    if err != nil {
        t.Errorf(err.Error())
        return
    }

    if currentval + 3 != currentnew {
    	t.Errorf("invalid new increment value %d want %d", currentnew, currentval + 3)
    	return
    }
}

// Тестирование RPC функции Api.SetSettings
func TestApiIncrementNumberMax(t *testing.T) {
    var err error
    var i uint

    // т.к. по заданию у нас нет возможности выставить текущее значение
    // счетчика через API, то используем такой workaround
    err = inc.SetSettings(1, 1)
    
    if err != nil {
        panic(err)
    }

    err = inc.SetSettings(30, 3)

    if err != nil {
        panic(err)
    }

    // ===========================================
    // инкрементация счетчика на 3
    // ===========================================

    for i = 1; i <= 13; i++ {

        if err := ApiIncrementNumber(i); err != nil {
            t.Errorf(err.Error())
            return
        }
    }

    // ===========================================
    // получение новго значения счетчика
    // ===========================================

    currentval, err := ApiGetNumber(311)

    if err != nil {
        t.Errorf(err.Error())
        return
    }

    if currentval != 6 {
        t.Errorf("invalid new increment value %d want 6", currentval)
        return
    }

}

// Функция вызова Api.IncrementNumber
func ApiIncrementNumber(id uint) error {
    var rpcres jsonRpcResult

    recorder := httptest.NewRecorder()

    reqincnum, err := newPostRpcJsonReq(id, "Api.IncrementNumber", cfg.App.ApiPath, nil)

    if err != nil {
        return err
    }

    rpc.ServeHTTP(recorder, reqincnum)

    if status := recorder.Code; status != http.StatusOK {
        return fmt.Errorf("api wrong status code: %d want %d",
            status, http.StatusOK)
    }

    if err := json.Unmarshal([]byte(recorder.Body.String()), &rpcres); err != nil {
        return fmt.Errorf("api unexpected json rpc body: %s", err.Error())
    }

    if rpcres.Id != id {
        return fmt.Errorf("rpc api unexpected ID %d want %d", rpcres.Id, id)
    }

    if rpcres.Error != nil {
        return fmt.Errorf("rpc api error: %v", rpcres.Error)
    }

    return nil
}

// Функция вызова Api.GetNumber
func ApiGetNumber(id uint) (int, error) {
    var rpcres jsonRpcResult

    recorder := httptest.NewRecorder()

    reqgetnew, err := newPostRpcJsonReq(id, "Api.GetNumber", cfg.App.ApiPath, nil)

    if err != nil {
        return 0, err
    }

    rpc.ServeHTTP(recorder, reqgetnew)

    if status := recorder.Code; status != http.StatusOK {
        return 0, fmt.Errorf("api wrong status code: %d want %d",
            status, http.StatusOK)
    }

    if err := json.Unmarshal([]byte(recorder.Body.String()), &rpcres); err != nil {
        return 0, fmt.Errorf("api unexpected json rpc body: %s", err.Error())
    }

    if rpcres.Id != id {
        return 0, fmt.Errorf("rpc api unexpected ID %d want %d", rpcres.Id, id)
    }

    if rpcres.Error != nil {
        return 0, fmt.Errorf("rpc api error: %v", rpcres.Error)
    }

    if !isIfaceNumber(rpcres.Result) {
        return 0, fmt.Errorf("rpc api invalid result type, want number: %v", rpcres.Result)
    }

    return int(rpcres.Result.(float64)), nil
}

// Функция вызова Api.SetSettings
func ApiSetSettings(id uint, max_val, inc_val int) error {
    var rpcres jsonRpcResult

    recorder := httptest.NewRecorder()

    reqsetstng, err := newPostRpcJsonReq(id,
                        "Api.SetSettings", 
                        cfg.App.ApiPath, map[string]int{
                            "max_val": max_val,
                            "inc_val": inc_val,
                        })

    if err != nil {
        return err
    }

    rpc.ServeHTTP(recorder, reqsetstng)

    if status := recorder.Code; status != http.StatusOK {
        return fmt.Errorf("api wrong status code: %d want %d",
            status, http.StatusOK)
    }

    if err := json.Unmarshal([]byte(recorder.Body.String()), &rpcres); err != nil {
        return fmt.Errorf("api unexpected json rpc body: %s", err.Error())
    }

    if rpcres.Id != id {
        return fmt.Errorf("rpc api unexpected ID %d want %d", rpcres.Id, id)
    }

    if rpcres.Error != nil {
        return fmt.Errorf("rpc api error: %v", rpcres.Error)
    }

    return nil
}


// функция для проверки типа interface{}
// используется для проверки типа результата RPC запроса
func isIfaceNumber(val interface{}) bool {
	switch val.(type) {
	case float64:
		return true
	}
	return false
}

// функция создает JSON RPC 2.0 запрос
func newPostRpcJsonReq(id uint, method, path string, params interface{}) (*http.Request, error) {

	jsonparams := jsonRpcRequest {
		Jsonrpc: "2.0",
		Method:  method,
		Params:  []interface{}{},
		Id:      id,
	}

	if params != nil {
		jsonparams.Params = []interface{}{ params }
	}

	data, err := json.Marshal(jsonparams)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", path, bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")

	return req, nil
}