---
description: API Design Guidelines for Brayat
---

# API Design Rules for Brayat

These rules define the RESTful API design standards for the Brayat application backend.

## 1. General Principles

- **Keep it Simple**: Use standard HTTP methods correctly (GET, POST, PUT, DELETE).
- **JSON only**: Request and response bodies must always be `application/json` (except for file uploads).
- **Stateless**: The API remains stateless, using session short-codes or link-codes for authorization.
- **Versioned base path**: Group all API routes under `/api/v1/`.

## 2. URL Structures

Resource URIs should be nouns, plural, and lowercase.
Use nested routes only if the relationship is strictly hierarchical.

**Good:**

- `GET /api/v1/sessions/{sessionID}`
- `POST /api/v1/sessions`
- `GET /api/v1/sessions/{sessionID}/people`
- `POST /api/v1/sessions/{sessionID}/people`
- `DELETE /api/v1/people/{personID}`

## 3. Standard Request/Response Formats

### Success Response

Always return the created or updated object on POST/PUT requests wrapped in a `data` object.

```json
{
  "data": {
    "id": "123",
    "name": "Budi"
  }
}
```

### Error Response

Always return errors in a consistent outer format with appropriate HTTP status codes.

```json
{
  "error": {
    "code": "invalid_input",
    "message": "Name is required."
  }
}
```

## 4. Polling & Syncing

Since we use a polling model (~5s), the API must be heavily optimized for fast reads. Consider `If-None-Match` or `last_updated_at` query params.

## 5. Security & Validation

- Validate link codes (admin or viewer) passed via `Authorization: Bearer <link_code>`.
- Always sanitize and validate client inputs.
