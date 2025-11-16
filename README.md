# Тестовое задание Авито
## Описание проекта
### Использованные технологии:
<img width="110" height="110" title="Golang" alt="golang" src="https://github.com/user-attachments/assets/9218acc9-6b90-4dbd-9ef4-4403f9664312"/>
<img width="100" height="100" title="Gin" alt="gin" src="https://github.com/user-attachments/assets/907271ec-aa81-4202-b9c7-5c945e1c4abb" />
<img width="110" height="100" title="Gorm" alt="gorm" src="https://github.com/user-attachments/assets/4c6cc1f2-1d94-47ca-b4af-8bea9f1becc0" />
<img width="110" height="110" title="PostgreSQL" alt="postgresql" src="https://github.com/user-attachments/assets/82b28ccc-f900-4877-a3d4-c19b439e180d"/>
<img width="110" height="110" title="Docker" alt="docker" src="https://github.com/user-attachments/assets/3b3e98c7-2d7f-4013-b603-f5694bd174e9" />

### Реализованный функционал:
- Создание команды с участниками
- Получение информации об определенной команде
- Возможность установить статус is_active определенному участнику
- Получение списка PR, на которые участник назначен в качестве ревьюера
- Создание PR с автоматическим назначением 2-х случайных участников команды
- Merge PR
- Переназначение определнного ревьюера на PR
- Просмотр статистики кол-ва PR, на которые назначены участники
- Массовая деактивация участников определенной команды
- Переназначение assigned_reviewers у всех PR определенной команды

### Установка и запуск (Без использования Docker)
1. Склонируйте репозиторий
```
git clone https://github.com/Alexeyts0Y/TEST_TASK_AVITO
cd TEST_TASK_AVITO
```

2. Создайте файл .env по шаблону .env.template
```
# .env

# Имя вашей базы данных
DB_NAME=your_db_name

# Хост вашей базы данных
DB_HOST=your_db_host

# Пароль от вашей базы данных
DB_PASSWORD=your_db_password

# Порт вашей базы данных
DB_PORT=your_db_port

# Имя пользователя базы данных
DB_USER=your_db_user

# Порт на котором будет работать сервер
SERVER_PORT=your_server_port
```

3. Запустите Makefile скрипт
```
make run
```

### Установка и запуск (С использованием Docker)
1. Склонируйте репозиторий
```
git clone https://github.com/Alexeyts0Y/TEST_TASK_AVITO
cd TEST_TASK_AVITO
```

2. Создайте файл .env по шаблону .env.template
```
# .env

# Имя вашей базы данных
DB_NAME=your_db_name

# Хост вашей базы данных
DB_HOST=your_db_host

# Пароль от вашей базы данных
DB_PASSWORD=your_db_password

# Порт вашей базы данных
DB_PORT=your_db_port

# Имя пользователя базы данных
DB_USER=your_db_user

# Порт на котором будет работать сервер
SERVER_PORT=your_server_port
```

3. Запустите Makefile скрипт
```
make docker-run
```

### Структура проекта
```
TEST_TASK_AVITO/
├── api/
│   └── openapi.yaml                    # OpenAPI спецификация
├── cmd/
│   └── server/
│       └── main.go                     # Точка входа в приложение
├── internal/
│   ├── config/
│   │   └── config.go                   # Конфигурация приложения
│   ├── errors/
│   │   └── errors.go                   # Кастомные ошибки
│   ├── handler/                        # Хендлеры (разделены по доменам)
│   │   ├── handler.go                  # Базовая структура Server и общие функции
│   │   ├── team_handlers.go            # Хендлеры для команд
│   │   ├── user_handlers.go            # Хендлеры для пользователей
│   │   ├── pull_request_handlers.go    # Хендлеры для пул-реквестов
│   │   └── stats_handlers.go           # Хендлеры для статистики
│   ├── model/                          # Модели данных (разделены)
│   │   ├── model.go                    # Базовая модель (BaseModel)
│   │   ├── user.go                     # Модель User и методы
│   │   ├── team.go                     # Модель Team
│   │   └── pull_request.go             # Модель PullRequest и методы
│   ├── repository/                     # Репозитории (разделены по доменам)
│   │   ├── repository.go               # Основной интерфейс репозитория и структура
│   │   ├── team_repository.go          # Репозиторий для команд
│   │   ├── user_repository.go          # Репозиторий для пользователей
│   │   ├── pull_request_repository.go  # Репозиторий для пул-реквестов
│   │   └── stats_repository.go         # Репозиторий для статистики
│   └── utils/
│       └── choose_random_candidates.go # Утилита для выбора случайных кандидатов
├── pkg/
│   └── api/
│       └── api.gen.go                  # Сгенерированный код из OpenAPI
├── .env.template                       # Шаблон для переменных окружения
├── .gitignore
├── .golang-ci.yml                      # Конфигурация golangci-lint
├── go.mod
├── go.sum
├── oapi-codegen.yaml                   # Конфигурация для oapi-codegen
└── README.md
```
