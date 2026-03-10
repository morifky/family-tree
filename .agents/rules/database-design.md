---
description: Database Design & SQLite Guidelines for Brayat
---

# Database Design Rules (SQLite)

These rules define the database standards and configuration for Brayat using SQLite.

## 1. Connection Settings & Pragmas

SQLite must be optimized for concurrent reads and safe writes. Every connection should execute:

```sql
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA foreign_keys = ON;
PRAGMA busy_timeout = 5000;
```

## 2. Table Design & IDs

- Use `TEXT` for primary keys (nanoid or short UUID).
- Avoid sensitive data (no birthdates, addresses, etc.).
- Always include `created_at` and `updated_at`.

## 3. Relationships & Constraints

- Explicitly define `FOREIGN KEY` constraints.
- Use `ON DELETE CASCADE` where appropriate to prevent orphaned records.

## 4. Query Patterns

- Use prepared statements to prevent SQL injection.
- Keep read queries flat and simple (`SELECT * FROM people WHERE session_id = ?`). Graph assembly happens on the frontend.
- Wrap multi-table inserts/updates in transactions.

## 5. Migrations

- Keep `up.sql` and `down.sql` for simple schema migrations. apply on startup.
