-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- таблица хранения значений `Incrementor`
CREATE TABLE IF NOT EXISTS `incrementor`
(
    `id`             INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    `key`            VARCHAR(128)    NOT NULL,
    `val`            BIGINT UNSIGNED NOT NULL,
    `created_at`     TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`     TIMESTAMP       NOT NULL 
                                     DEFAULT CURRENT_TIMESTAMP
                                     ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);

-- уникальный индекс для ключевого поля
-- для предотвращения дублирования ключей
CREATE UNIQUE INDEX `incrementor_key`
    ON `incrementor` (`key`);

-- значения "по умолчанию"
INSERT INTO `incrementor` (`key`, `val`)
    VALUES
        ("max_val", 0), -- максимальное значение
        ("cur_val", 0), -- текущее значение
        ("inc_val", 1); -- значение приращения

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

-- удаление индекса
ALTER TABLE `incrementor` DROP INDEX `incrementor_key`;

-- удаление таблицы
DROP TABLE IF EXISTS `incrementor`;