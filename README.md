# EateryGlass
## Система бронирования столиков в ресторане

[Примеры использование API в Postman](https://documenter.getpostman.com/view/21404316/Uz5Nishh)

[Документация](/docs/README.md)

---

## Установка и запуск

### 1. Клонируем репозиторий

### 2. Необходимо запустить отдельный контейнер с PostgresQL. 
Например

```bash
$ docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
```

### 3. Дальше нужно сконфигурировать под себя `.env` файл

```bash
# database

# IP контейнера с PostgresQL (см. ниже) 
DB_HOST_ADDR=172.17.0.2 
# Порт на котором в контейнера запущен PostgresQL 
# (по умолчанию 5432)
DB_HOST_PORT=5432
# Название базы данных с которой собираемся работать
# (по умолчанию  postgres)
DB_NAME=postgres

# Имя пользователя в PostgresQL, через которого
# собираемся взаимодействовать с БД
DB_USERNAME=postgres
# Пароль к этому пользователю
# (POSTGRES_PASSWORD который указывали в пункте 2)
DB_PASSWORD=mysecretpassword

# app

APP_HOST=0.0.0.0
# Порт который будет слушать контейнер с приложением
APP_PORT=8000

# dir
# Рабочая директория внутри контейнера с приложением
# (менять не нужно)
MAIN_DIR=/usr/src/app
```

#### __Как получить IP адрес контейнера с PostgresQL__

Нужно получить адреса контейнеров в docker-сети

```bash
$ docker network inspect bridge 
```
```json
[
    {
        "Name": "bridge",

        .........

        "Containers": {
            "16fac93.........8ef7b5": {
                "Name": "some-postgres",
                "EndpointID": "c51e505a38ca83d2f206d7f0a7e6990db0f4d926e3b81bf42b727c64c642ea1c",
                "MacAddress": "02:42:ac:11:00:02",
                "IPv4Address": "172.17.0.2/16",
                "IPv6Address": ""
            }
        }

        .........

    }
]
```

Нужно найти контейнер с указанным выше именем, в поле `IPv4Address` и будет нужный IP адрес.

### 4. Собираем docker-образ приложения

В корневой директории приложения выполните

```bash
$ docker build -t eateryglass .
```

Эта команда соберет docker-образ приложения с именем `eateryglass` (имя можно менять или оставить пустым) 

### 5. Запустим приложение

В отдельном терминале выполните команду

```bash
$ docker run -it --name EateryGlass eateryglass
```

Команда запустит контейнер `EateryGlass` из образа `eateryglass` (имя можно поменять или оставить пустым)


Дальше приложение само инициализирует базу данных, создаст необходимые таблицы и заполнит их тестовыми данными.

Теперь осталось узнать IP адрес контейнера, в котором запущено приложение (_например_ методом описанным выше) и отправлять запросы.

[Примеры использование API в Postman](https://documenter.getpostman.com/view/21404316/Uz5Nishh)

[Документация](/docs/README.md)
