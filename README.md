
#  Market Rent Bot

Телеграм-бот для управления арендой павильонов и арендаторами на рынке, написанный на Go с использованием `telego` и SQLite.

## 🚀 Функционал

- **Для админа**
  - Управление пользователями: добавление, фильтрация, удаление.
  - Управление павильонами: просмотр, добавление, поиск, постраничный вывод.
  - Управление арендаторами: многоэтапное добавление с выбором видов деятельности, договорами и связью с павильоном.
  - Управление видами деятельности: просмотр, добавление, множественный выбор.
  - Импорт данных из CSV/Excel (заготовка).

- **Для арендаторов**
  - Просмотр платежей.
  - Просмотр своего договора.

## 🛠️ Установка

1. Клонируйте репозиторий и войдите в папку:
   ```bash
   git clone https://github.com/AlexEagle1535/market-rent-bot.git
   cd market-rent-bot
   ```

2. Установите зависимости:
   ```bash
   go mod download
   ```

3. Укажите токен бота и юзернеймы администратров в файле .env:
   ```bash
   BOT_TOKEN=YOUR_TOKEN
   ADMINS=user1,user2,user3
   ```

4. Запустите бота:
   ```bash
   go run main.go
   ```

5. При первом запуске создаётся `market.db` с необходимыми таблицами.

## ⚙️ Структура

- `db/` — работа с SQLite: инициализация базы, CRUD для пользователей, арендаторов, павильонов, договоров, видов деятельности.
- `handlers/` — обработчики команд и callback-запросов.
- `menu/` — функции генерации inline-клавиатур.
- `states/` — FSM-состояния пользователей и временное хранилище введенных данных.
- `main.go` — точка входа, подключение Telegram-бота и маршрутизация.

## 🧩 Используемые библиотеки

- [`github.com/mymmrac/telego`](https://pkg.go.dev/github.com/mymmrac/telego)
- `modernc.org/sqlite` — SQLite-драйвер.

## 🔧 Архитектура

- **FSM-состояния** хранятся в `states`: строки (`userStates`), список (`ListState`), временные данные (`tempStorage`).
- **Обработчики**:
  - `TextMessage`: обрабатывает ввод текста по состояниям.
  - `CallbackQuery`: обрабатывает inline-кнопки, переключение меню.
- **Меню**: `menu/*.go` — функции для генерации клавиатур.
- **База данных**: `db/*.go` — функции CRUD, `InitDB()` создаёт таблицы.

## 📝 Особенности

- Поддержка **многоэтапного ввода** (например, добавление арендатора, павильона, договора).
- **Постраничный вывод** списков с `ListState.Scope` и `Page`.
- **Множественный выбор** видов деятельности с накоплением через `select_activity_type:` и кнопку `finish_activity_selection`.
- **Связи** между таблицами (`users ↔ tenants`, `tenants ↔ activity_types`, `tenants ↔ pavilions`) с внешними ключами.
