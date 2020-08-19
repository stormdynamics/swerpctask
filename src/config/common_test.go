// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package config

import (
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
)

// Тестирование функции формирования адреса вида <хост>:<port>
// значения по умолчанию
func TestAppHostPortDefault(t *testing.T) {
	var   val  string = ""
	const want string = "127.0.0.1:8080"
	
	cfg := NewConfig()

	if cfg == nil {
		t.Error("cannot create an instance of Config")
		return
	}

	if val = cfg.AppHostPort(); val != want {
		t.Errorf("invalid value %s want %s", val, want)
	}
}

// Тестирование функции формирования адреса вида <хост>:<port>
// с установкой значения через переменные окружения
func TestAppHostPort(t *testing.T) {
	var val string
	var cfg *Config

	testset := []struct {
		host string
		port string
		want string
	} {
		{ "",                "",     "127.0.0.1:8080"       },
		{ "",                "1288", "127.0.0.1:1288"       },
		{ "8.8.8.8",         "",     "8.8.8.8:8080"         },
		{ "0.0.0.0",         "7733", "0.0.0.0:7733"         },
		{ "192.168.0.1",     "2200", "192.168.0.1:2200"     },
		{ "10.20.30.50",     "6070", "10.20.30.50:6070"     },
		{ "133.221.222.231", "7132", "133.221.222.231:7132" },
	}

	for _, testval := range testset {
		os.Setenv("APP_HOST", testval.host)
		os.Setenv("APP_PORT", testval.port)

		cfg = NewConfig()

		if cfg == nil {
			t.Error("cannot create an instance of Config")
			continue
		}

		if val = cfg.AppHostPort(); val != testval.want {
			t.Errorf("invalid value %s want %s", val, testval.want)
		}
	}

	os.Unsetenv("APP_HOST")
	os.Unsetenv("APP_PORT")
}

// Тестирование функции формирования адреса БД вида <хост>:<port>
// значения по умолчанию
func TestDbHostPortDefault(t *testing.T) {
	var   val  string = ""
	const want string = "127.0.0.1:3306"
	
	cfg := NewConfig()

	if cfg == nil {
		t.Error("cannot create an instance of Config")
		return
	}

	if val = cfg.DbHostPort(); val != want {
		t.Errorf("invalid value %s want %s", val, want)
	}
}

// Тестирование функции формирования DSN строки для
// подключения к БД по умолчанию
func TestDbDSNDefault(t *testing.T) {

	cfg := NewConfig()

	if cfg == nil {
		t.Error("cannot create an instance of Config")
		return
	}

	dsncfg, err := mysql.ParseDSN(cfg.DbDSN())

	if err != nil {
		t.Error(err.Error())
	}

	if addr := cfg.DbHostPort(); dsncfg.Addr != addr {
		t.Errorf("invalid `Addr` value %s want %s", dsncfg.Addr, addr)
		return
	}

	if dsncfg.User != cfg.Db.User {
		t.Errorf("invalid `User` value %s want %s", dsncfg.User, cfg.Db.User)
		return
	}

	if dsncfg.Passwd != cfg.Db.Pass {
		t.Errorf("invalid `Passwd` value %s want %s", dsncfg.Passwd, cfg.Db.Pass)
		return
	}

	if dsncfg.DBName != cfg.Db.Name {
		t.Errorf("invalid `DBName` value %s want %s", dsncfg.DBName, cfg.Db.Name)
		return
	}

	if dsncfg.Collation != cfg.Db.Colt {
		t.Errorf("invalid `Collation` value %s want %s", dsncfg.Collation, cfg.Db.Colt)
	}
}

// Тестирование получения значения пути API для HTTP сервера по умолчанию
func TestAppApiPathDefault(t *testing.T) {
	var   val  string = ""
	const want string = "/api"
	
	cfg := NewConfig()

	if cfg == nil {
		t.Error("cannot create an instance of Config")
		return
	}

	if cfg.App.ApiPath != want {
		t.Errorf("invalid value %s want %s", val, want)
	}
}

// Тестирование получения значения пути API для HTTP сервера
func TestAppApiPath(t *testing.T) {
	var cfg *Config

	testset := []struct {
		path string
		want string
	} {
		{ "",          "/api"      },
		{ "/test/api", "/test/api" },
		{ "/app/_api", "/app/_api" },
	}

	for _, testval := range testset {
		os.Setenv("APP_API_PATH", testval.path)

		cfg = NewConfig()

		if cfg == nil {
			t.Error("cannot create an instance of Config")
			continue
		}

		if cfg.App.ApiPath != testval.want {
			t.Errorf("invalid value %s want %s", cfg.App.ApiPath, testval.want)
		}

	}

	os.Unsetenv("APP_API_PATH")
}