# Team Manager

Team Manager – это веб-приложение для эффективного управления ресурсами команды. Приложение
предлагает простой и интуитивно понятный UI, с малым порогом входа для новых пользователей.

## Стэк технологий

- **Go** (основной язык)
- **[Gin](https://github.com/gin-gonic/gin)** (HTTP-фреймворк)
- **gRPC** (для взаимодействия с [сервисом авторизации](https://github.com/hard-gainer/auth-service))
- **PostgreSQL** (в качестве основной базы данных, использовался драйвер [pgx](https://github.com/jackc/pgx))
- [Migrate](https://github.com/golang-migrate/migrate) (для миграций)
- **Tailwind CSS** (стилизация фронтенда)
- **htmx** (асинхронные операции без перезагрузки страницы)

## Основной функционал

- Регистрация и авторизация пользователей с использованием системы ролей (manager/employee) для создания и назначения задач. 
- Роль "manager" позволяет:
    - Создавать\удалять проекты
    - Создавать приглашения для добавления в проект новых участников
    - Регистририровать и назначать задачи в пределах проекта.
    - Просматривать статистику затраченного времени на задания внутри проектов.
    - Редактировать задачи внутри проекта.
- Роль "employee" позволяет:
    - Изменять статусы задач, например, принимая их в работу или завершая выполнение, при этом
      приложение производит расчет трудозатрат исполнителя.  

## Установка

1. Установите и настройте PostgreSQL.
2. Склонируйте репозиторий:
   ```
   git clone https://github.com/your-username/task-tracker.git
   ```
3. Установите зависимости:
   ```
   go mod tidy
   ```
4. Создайте нужные таблицы и запустите миграции:
   ```
   make createdb
   make migrateup 
   ```

## Использование

1. Запустите сервер:
   ```
   go run main.go
   ```
2. Откройте в браузере:
   ```
   http://localhost:8080
   ```
3. Перейдите на страницу регистрации или авторизации, чтобы начать работать с проектами и задачами.

## Особенности

- Микросервисная архитектура: аутентификация через внешний auth-сервис по протоколу gRPC.
- Поддержка асинхронных действий через htmx.
