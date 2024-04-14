# Сервис баннеров



## Установка и запуск

### Используя Docker

Для запуска **Сервис баннеров** требуются следующие шаги:

1. Рекомендуется выполнять запуск на Unix-систме (Ubuntu 22.04.3 LTS).
2. Установить [Docker и Docker-compose](https://www.docker.com/) на локальную машину.
3. Склонировать репозиторий проекта.
4. Собрать проект при помощи **Makefile:**
```bash
make up
```

### Без использования Docker

1. Рекомендуется выполнять запуск на Unix-систме (Ubuntu 22.04.3 LTS).
2. Должна быть установлена версия go до 1.22 (у меня go 1.21.6).
3. Неоходимо установить БД Postgres (у меня 16.1).
4. Выполнить файл baseline.sql, расположенный в папке ./internal/postgres/v1.0/v1.01.
5. Склонировать репозиторий проекта.
6. В файле main.go (./cmd/http/main.go) на 19 строчке необходимо в аргументе функции изменить название файла с
**.env** на **.env_localhost** (отличие у них в хосте):
было
```go
err := godotenv.Load(".env")
```
стало после замены
```go
err := godotenv.Load(".env_localhost")
```
7. Запустить сервер можно с помощью **Makefile:**:
```bash
make run
```

## База Данных и Кэш

БД состоит из двух таблиц: **banners** и **banners_data**. Таблица **banners** служит для быстрого поиска необходимого
id баннера, а **banners_data** - хранит непосредственно информацию банера.
БД удовлетворяет интерфесу **DBmodel** (см. ./internal/postgres/db_interface.go):
```go
type DBModel interface {
GetUserBanner(tagId, featureId, role int) (*dto.GetBanner, error)
GetBanners(whereStmt, limitOffsetStmt string) ([]dto.GetBanner, error)
InsertBanner(banner dto.PostPatchBanner) (int, error)
UpdateBannerId(id int, banner dto.PostPatchBanner) (*storage.UpdateDeleteFromDB, error)
DeleteBanner(id int) (*storage.UpdateDeleteFromDB, error)
}
```
Схема БД представлена на рисунке ниже:

Кэш удовлетворяет интерфейсу **Storage** (см. ./internal/storage/storage.go):
```go
type Storage interface {
Save(bannerId int, banner *dto.PostPatchBanner) error
Get(tagId, featureId int) (*BannerInfo, error)
Update(banner *UpdateDeleteFromDB) error
Delete(banner *UpdateDeleteFromDB) error
}
```
Считал, что получить информуцию из кэша можно только по пути ```/user_banner```. При запросах поостальным путям
информация в кэше **обновляется**.


## :bulb: API приложения

- ```GET /user_banner``` - получение баннера для пользователя и/или админа по tag_id. Пример получения списка продуктов:
```json
{
    "all_count": 33,
    "goods": [
        {
            "name": "milk",
            "size": "0.2m",
            "product_id": 1,
            "count": 15
        },
        {
            "name": "doors",
            "size": "1.8m",
            "product_id": 3,
            "count": 18
        }
    ]
}
```

- ```POST /product``` - добавление продуктов в таблицу продуктов. Пример добавления продукта:
```json
{
    "name": "milk",
    "size": "0.2m",
    "count": 15
}
```
- ```DELETE /product``` - удаление продукта из списка продуктов и из резерва. Пример работы удаления по следующему ID `[2]`:
```json
[
    {
        "name": "tables",
        "size": "1.5m",
        "product_id": 2,
        "count": 20
    }
]
```

- ```POST /product/warehouse``` - резервирование продуктов для дальнейшей доставки. На вход подается массив уникальных ID продуктов.
  Пример резервирования продуктов с ID `[1,2]`:
```json
[
    {
        "name": "milk",
        "size": "0.2m",
        "product_id": 1,
        "count": 15
    },
    {
        "name": "tables",
        "size": "1.5m",
        "product_id": 2,
        "count": 20
    }
]
```
- ```DELETE /product/warehouse``` - удаление продуктов из резерва. На вход подается массив уникальных ID продуктов.
  Пример удаления продуктов из резерва с ID `[1]`:
```json
[
    {
        "name": "milk",
        "size": "0.2m",
        "product_id": 1,
        "count": 15
    }
]
```

## Стартовые данные базы данных:

`Таблица products`

|  unique_code  |  product_name  |  size   |  count  |
|:-------------:|:--------------:|:-------:|:-------:|
|       1       |      milk      |  0.2m   |   15    |
|       2       |     tables     |  1.5m   |   20    |
|       3       |     doors      |  1.8m   |   18    |
|       4       |  microphones   |  0.15m  |    4    |

`Таблица warehouse`

| product_name | can_be_use |
| :----------: | :--------- |
|             |            |


## Тесты
`Для использования тестов через postman`
1. Необходимо скачать коллекцию [postman-тестов](/postman_collection)
2. Запустить тесты из коллекции

`Для использования unit-тестов`
```bash
go test ./
```
Требования к тестам см. в [postman-test.md](/backend/doc/postman-test.md) и [unit-test.md](/backend/doc/unit-test.md)