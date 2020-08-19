// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package config

import (
	"os"
	"testing"
)

// Тестирование создания экземпляра Config
func TestCreating(t *testing.T) {
	cfg := NewConfig()

	if cfg == nil {
		t.Error("cannot create an instance of Config")
	}
}

// Тестирование функции getEnvString с получением значения по умолчанию
func TestGetEnvStringDefault(t *testing.T) {
	const key string = "CONFIG_TEST_ENV"
	const val string = "ENV-TEST-DEFAULT-STRING"

	if enval := getEnvString(key, val); enval != val {
		t.Errorf("invalid environment variable value `%s` want `%s`", enval, val)
	}
}

// Тестирование функции getEnvString с получением значения из
// переменной окружения
func TestGetEnvString(t *testing.T) {
	const key string = "CONFIG_TEST_ENV"
	const val string = "ENV-TEST-STRING"

	os.Setenv(key, val)

	if enval := getEnvString(key, ""); enval != val {
		t.Errorf("invalid environment variable value `%s` want `%s`", enval, val)
	}

	os.Unsetenv(key)
}

// Тестирование функции getEnvUint16 с получением значения по умолчанию
func TestGetEnvUint16Default(t *testing.T) {
	const key string = "CONFIG_U16DEF_TEST_ENV"
	const val uint16 = 7331

	if enval := getEnvUint16(key, val); enval != val {
		t.Errorf("invalid environment variable value %d want %d", enval, val)
	}
}

// Тестирование функции getEnvUint16 с получением значения из
// переменной окружения
func TestGetEnvUint16(t *testing.T) {
	const key string = "CONFIG_U16_TEST_ENV"

	testset := []struct {
		val  string
		want uint16
	} {
		{ "1337",   1337   },
		{ "6124",   6124   },
		{ "5",      5      },
		{ "55555",  55555  },
	}

	for _, testval := range testset {
		os.Setenv(key, testval.val )

		if enval := getEnvUint16(key, 0); enval != testval.want {
			t.Errorf("invalid environment variable value %d want %d", enval, testval.want)
		}
	}
	
	os.Unsetenv(key)
}