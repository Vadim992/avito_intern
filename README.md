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
Для запуска используется порт ```:3000```
### ВАЖНО
При запуске автоматически запоняю БД (в **banners** миллион элементов, в **banners_data** - 1000).
Чтобы убрать автозаполнение БД необходимо закомментировать следующий код в main.go (./cmd/http/main.go) на 52 строчке:
```go
	if err := DB.FillDb(); err != nil {
logger.ErrLog.Fatalf("cannot fill DB: %v", err)
}
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
![drawSQL-image-export-2024-04-14](https://github.com/Vadim992/avito_intern/assets/105795306/df37ccdb-f9c8-4275-aed0-93cef914396a)

Кэш удовлетворяет интерфейсу **Storage** (см. ./internal/storage/storage.go):
```go
type Storage interface {
Save(bannerId int, banner *dto.PostPatchBanner) error
Get(tagId, featureId int) (*BannerInfo, error)
Update(banner *UpdateDeleteFromDB) error
Delete(banner *UpdateDeleteFromDB) error
}
```
Структра кэша - **InMemoryStorage** (см. ./internal/storage/storage.go):
```go
type InMemoryStorage struct {
SearchStorage map[SearchIds]int
Banners       map[int]*BannerInfo
mu            sync.RWMutex
}
type SearchIds struct {
TagId     int
FeatureId int
}

type BannerInfo struct {
BannerId  int
Content   dto.BannerContent
IsActive  bool
UpdatedAt time.Time
}
```

Первая мапа ```SearchStorage map[SearchIds]int``` соответсвует таблице ```banners``` в БД, а вторая
```Banners       map[int]*BannerInfo``` - ```banners_data``` (без поля created_at, так как оно не нужно обычным
пользователям).
Считал, что получить информуцию из кэша можно только по пути ```/user_banner```. При запросах поостальным путям
информация в кэше **обновляется**.


## API приложения

-```Header 'token'``` - токены админа и пользователя храню в **.env** файле. Так как по условию не было сказано,
что их роль пользователя передает front d jwt токене, то для упрощения принял такое решение.

- ```GET /user_banner``` - получение баннера для пользователя и/или админа по tag_id и feature_id. Если данные в кэше
  не устарели (не прошло 5 мин с момента последнего обновления кэша), то получаю их оттуда (при условии, что
  ```use_lasr_revision=false```). Для успешного ответа используется структура (см. ./internal/dto/dto.go):
```go
type BannerContent struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
	Url   *string `json:"url"`
}
```

Для ответа кода **400** используется структура  (см. ./internal/req_err_handle.go):
```go
type ErrorStruct struct {
Err string
}
```

- ```GET /banner``` - получение требуемых данных от админа, считал что админу нужна только актуальная информация из БД,
  поэтому при этом запросе кэш **не** учавствует. Структуры для успешного ответа представлены ниже
  (см. ./internal/dto/dto.go):
```go
type BannerContent struct {
  Title *string `json:"title"`
  Text  *string `json:"text"`
  Url   *string `json:"url"`
}
type PostPatchBanner struct {
  FeatureId *int           `json:"feature_id"`
  TagIds    []int64        `json:"tag_ids"`
Content   *BannerContent `json:"content"`
IsActive  *bool          `json:"is_active"`
}

type GetBanner struct {
BannerId *int `json:"banner_id"`
PostPatchBanner
CreatedAt *time.Time `json:"created_at"`
UpdatedAt *time.Time `json:"updated_at"`
}
```
Ответ приходит массивом структур ```GetBanner```.


- ```POST /banner``` - занесение нового банера в таблицу, результат ответа json структуры (см. ./internal/dto/dto.go):
```go
type BannerId struct {
BannerId int `json:"banner_id"`
}
```
После занесения в БД **обновляю кэш**.
### Примечание
Если в массиве tag_ids были повторяющиеся тэги, то вставка проходит успешно, если хотя бы один из тэгов создает
конфликтную ситуацию при втсавке в БД - вставка отменяется и кэш не обновляется (код **400**).


- ```PATCH /banner/{id}``` - изменение данных, ситуация с tag_ids аналогична описанной выше. При обновлении tag_id в
  таблице ```banners``` старые tag_id удаляются и заменяются на полностью новые. После обновления данных в БД происходит
  обновление кэша. Метод ```PATCH``` возвращает структуру (см. ./internal/storage/storage.go):
```go
type UpdateDeleteFromDB struct {
  BannerId      int
  FeatureId     *int
  TagIds        []int64
ReqBanner     *dto.PostPatchBanner
UpdatedBanner *dto.PostPatchBanner
}
```
Она сипользуется для **обновления** данных в кэше (пользователь о ней ничего не знает).
- ```DELETE /banner/{id}``` - удаление банера по id. Метод ```DELETE``` возвращает такую структуру как и метод
- ```PATCH``` (пользователь об это также ничего не знает), которая служит для удаления данных их кэша.

## Тесты
`Для использования тестов через postman`
1. Необходимо скачать коллекцию **./postman**
2. Запустить тесты из коллекции

`Для использования unit-тестов`
Запуск тестов
```bash
make test
```
Был оттестирован ```GET /user_banner```. Тест расположен в файле ./internal/request_test.go. Для теста создавался mock
БД (см. ./internal/mockDB.go)
