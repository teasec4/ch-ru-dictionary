# БКРС Словарь - Backend

## Обзор

REST API для китайско-русского словаря БКРС. База данных содержит ~2M записей.

## База данных

**Файл:** `dictionary.db` (SQLite, ~700MB)

### Таблицы

```sql
CREATE TABLE entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hanzi TEXT NOT NULL UNIQUE,
    pinyin TEXT,
    pinyin_normalized TEXT
);

CREATE TABLE meanings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entry_id INTEGER NOT NULL REFERENCES entries(id),
    text TEXT NOT NULL,
    part_of_speech TEXT,
    order_num INTEGER DEFAULT 0
);

CREATE TABLE examples (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    meaning_id INTEGER REFERENCES meanings(id),
    chinese TEXT NOT NULL,
    translation TEXT
);

CREATE TABLE refs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    meaning_id INTEGER REFERENCES meanings(id),
    target_entry_id INTEGER REFERENCES entries(id),
    target_hanzi TEXT
);
```

### Индексы

```sql
CREATE INDEX idx_entries_hanzi ON entries(hanzi);
CREATE INDEX idx_entries_pinyin_norm ON entries(pinyin_normalized);
CREATE INDEX idx_meanings_entry ON meanings(entry_id);
```

### Статистика

| Таблица | Записей |
|---------|---------|
| entries | 2,081,762 |
| meanings | 1,986,420 |

## Структура проекта

```
backend/
├── cmd/
│   ├── api/main.go          # HTTP сервер (Gin)
│   └── cli/main.go          # CLI клиент для тестирования
├── internal/
│   ├── handler/
│   │   └── handler.go       # HTTP обработчики
│   └── search/              # Поиск (разбит по типам)
│       ├── hanzi.go        # Поиск по иероглифам
│       ├── pinyin.go       # Поиск по пиньинь
│       └── meaning.go      # Поиск по значению + автодополнение
├── go.mod
└── dictionary.db
```

## API Endpoints

| Endpoint | Описание |
|----------|-----------|
| `GET /api/entries?hanzi=北京` | Поиск по иероглифам |
| `GET /api/entries?pinyin=bei` | Поиск по пиньинь |
| `GET /api/entries?q=Пекин` | Поиск по русскому значению |
| `GET /api/entries?hanzi=北京&page=2&limit=20` | Пагинация |
| `GET /api/autocomplete?prefix=北京&limit=10` | Автодополнение |

## Пример ответа

```json
{
  "data": [
    {
      "hanzi": "北京",
      "pinyin": "běijīng",
      "meanings": [
        {"index": 1, "text": "1) Пекин ( город и административный район, КНР )"}
      ]
    }
  ],
  "total": 206,
  "page": 1,
  "limit": 20
}
```

## Запуск

```bash
# Запуск сервера
go run ./cmd/api

# Тестирование через CLI
go run ./cmd/cli/main.go 北京

# Или curl
curl "http://localhost:8080/api/entries?hanzi=北京"
```

## Логика поиска по иероглифам

1. **Точное совпадение** (`hanzi = '北京'`) - приоритет
2. **Префикс** (`hanzi LIKE '北京%'`) - слова начинающиеся на 北京
3. **Содержит** (`hanzi LIKE '%北京%'`) - слова содержащие 北京

## Известные проблемы

- FTS5 не работает (SQLite без FTS5 модуля) - используется LIKE

## TODO

- [ ] Rate limiting
- [ ] Кэширование (Redis)
- [ ] FTS5 для полнотекстового поиска