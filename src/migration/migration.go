// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

/*
    Встроенный менеджер управления миграциями на основе проекта GOOSE
    https://github.com/pressly/goose (документацию можно найти по ссылке)
*/

package migration

import (
	"os"
	"log"
	"fmt"
	"flag"
	"database/sql"

	"github.com/pressly/goose"
)

const (
    usageInfoPrefix = `Usage: migrate [OPTIONS] COMMAND
Examples:
    migrate status

Options:`

	usageInfoCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)

var (
	migrationFlags *flag.FlagSet
)

// Функция миграции, ожидает передачу через командную строку
// аргумента `migrate`.
// 
// ./app migrate выведет справочную информацию
// по командам менеджера миграций
func Migration(db *sql.DB) {

    if len(os.Args) < 2 || os.Args[1] != "migrate" {
        return
    }

	migrationFlags  = flag.NewFlagSet("migrate", flag.ExitOnError)
	migrationDir   := migrationFlags.String("dir", "./migrations", "directory with migration files")

	migrationFlags.Usage = printUsage
	migrationFlags.Parse(os.Args[2:])
	args := migrationFlags.Args()

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		migrationFlags.Usage()
		os.Exit(0)
	}

    // В данном приложении используется MySQL
    if err := goose.SetDialect("mysql"); err != nil {
        log.Fatal(err)
    }

	command := args[0]

	switch command {
    case "create":
        if err := goose.Run("create", nil, *migrationDir, args[1:]...); err != nil {
            log.Fatalf("migrate run error: %s", err.Error())
        }
    case "fix":
        if err := goose.Run("fix", nil, *migrationDir); err != nil {
            log.Fatalf("migrate run error: %s", err.Error())
        }
    default:
        if err := goose.Run(command, db, *migrationDir, args[1:]...); err != nil {
            log.Fatalf("migrate run error: %s", err.Error())
        }
    }

    os.Exit(0)
}

// Печатает справку
func printUsage() {
    fmt.Println(usageInfoPrefix)
    migrationFlags.PrintDefaults()
    fmt.Println(usageInfoCommands)
}