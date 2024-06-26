# av-banner-task
### Сервис управления баннерами
[Тестовое задание от Авито.](https://github.com/avito-tech/backend-trainee-assignment-2024)

#### Запуск и остановка сервиса в Docker'e:
- `make dc-up` запустит сервис и окружение
- `make dc-up-d` запустит сервис и окружение в демоне
- `make dc-stop` остановит контейнеры
- `make dc-down` остановит и удалит контейнеры

---

#### Приложение:
- Запущенное приложение будет доступно по адресу http://127.0.0.1:8080

---

### API методы:
Примеры запросов описанные для `Postman` в `docs/banner_service.postman_collection.json`

---

#### Swagger UI:
- Swagger описание будет доступно по адресу http://127.0.0.1:8080/swagger/index.html

---

### Что использовалось
- Код сервиса написан на `Golang` версии `1.22.0`
- HTTP роутинг на `chi` + `render` для обработки запросов и ответов
- Хранение баннеров в `PostgreSQL`, с покдлючением через `pgx/pgxpool`
- In-memory хранилищие `Redis`, с подключением через `go-redis`
- Конейниризация сервиса посредством `Docker`, c конфигурацией при помощи `docker-compose`
- [Документация к API](http://127.0.0.1:8080/swagger/index.html) написано при помощи `Swagger`.

---

### Детали реализации по заданию
Логика работы сервиса разделена на три уровня `handler`, `controller`, `storage`. 
В `handler` описаны API энтрипоинты и происходит обработка запроса и ответа.
В `controller` происходит основная обработка бизнес логики.
И в `storage` описывается взаимодействие с базами данных.
Также подключено логирование с использованием `logrus` и обработка параметров конфигурации сервиса из переменных окружения.


> Для авторизации доступов должны использоваться 2 вида токенов: пользовательский и админский.

Авторизация доступов реализована с использованием JWT токенов в которые зашифрована роль пользователя.
Для проверки написано `middleware` которое перехватывает запросы и валидирует токены.


> Реализуйте интеграционный или E2E-тест на сценарий получения баннера.

Реализован интеграционный тест для получения баннера. Для его прогона используется тестовое окружение, настроено в `docker-compose-api-test.yml`.

Команда запуска теста:
- `make integration-test`

> Если при получении баннера передан флаг use_last_revision, необходимо отдавать самую актуальную информацию. В ином случае допускается передача информации, которая была актуальна 5 минут назад.

Если флаг отсутствует или выставлен в `false`, то баннер достается из кэша хранилища `Redis` при его наличии, иначе происходит запрос к основному хранилищу в `PostgreSQL`. 
Срок жизни кэшированных баннеров регулируется переменной окружения `REDIS_EXPIRATION_DURATION`.
Кэшируются только баннеры которые запрашивал обычный пользователь.

> Баннеры могут быть временно выключены. Если баннер выключен, то обычные пользователи не должны его получать, при этом админы должны иметь к нему доступ.

Если пришел запрос с токеном от обычного пользователя, то ему отдаются только активные баннеры. При запросе с токеном от админа, в ответ придет запрашиваемый баннер игнорируя его состояние активности.

