# test-avito-mx

## Требования к оборудования для запуска
* Git
* Docker
* Docker Compose

## Инструкция по запуску проекта
1. Hеобходимо склонировать репозиторий
2. Запустить сервис из папки с проектом с помощью команды: docker-compose up

#### Для загрузки товаров нужно выполнить GET запрос на создание задания, например:

```
curl --location --request GET 'http://localhost:6000/api/v1/goods?url=http://simulator:8090/api/load-file&merchant-id=1'
```

где url - это ссылка на файл, merchant-id - id продавца, к чьему аккаунту будут привязаны загружаемые товары.

#### Для того чтобы узнать статус задания и краткую статистику по загрузке файла, нужно выполнить GET запрос, например:

```
curl --location --request GET 'http://localhost:6000/api/v1/progress?task-id=1'
```

где task-id - это id задания

#### Для проверки работоспособности сервиса нужно выполнить GET запрос, с помощью которого можно достать список товаров из базы, например:

```
curl --location --request GET 'http://localhost:6000/api/v1/list?merchant-id=1&offer-id=1' --data-urlencode 'name=
теле'
```

где merchant-id - это id продавца, offer-id - это уникальный идентификатор товара в системе продавца, name - это название товара
