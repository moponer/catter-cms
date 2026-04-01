---
name: catter-cms
description: Catter CMS content management integration. Manage posts, pages, categories, tags, media, and comments via REST API.
version: 1.0.0
---

# Catter CMS Skill

## When to use this skill

Use this skill when the user asks you to:

- Create, read, update, or delete **articles/posts**
- Manage **pages** (about, contact, etc.)
- Organize **categories** and **tags**
- Browse or manage **media files**
- Review or moderate **comments**
- Check or update **site settings**
- Any task involving CMS content management

## Authentication

All API requests require authentication via Bearer token.

1. Read the API key from environment variable `CATTER_CMS_KEY`
2. Read the CMS base URL from environment variable `CATTER_CMS_URL`
3. Include the header: `Authorization: Bearer <value of CATTER_CMS_KEY>`
4. The key starts with `cat_` prefix

**Important:** Never expose the raw API key in responses. Always reference it from the environment variable.

## API Base

All endpoints are relative to the CMS base URL stored in `CATTER_CMS_URL`.

Example: if `CATTER_CMS_URL` is `https://example.com`, then the full URL for listing posts is `https://example.com/api/v1/posts`.

## Request & Response Format

- **Content-Type:** application/json
- **Pagination:** `page` (default 1), `page_size` (default 10, max 100)
- **Search:** `keyword` query parameter for text search

**Response structure (all endpoints):**

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

`code = 0` means success. Any non-zero code indicates an error.

## API Endpoints

### Articles (Posts)

| Method | Endpoint | Description | Required Permission |
|--------|----------|-------------|-------------------|
| GET | /api/v1/posts | List posts (supports `page`, `page_size`, `keyword`, `status`, `category_id`) | read |
| GET | /api/v1/posts/:id | Get post detail | read |
| POST | /api/v1/posts | Create post | write |
| PUT | /api/v1/posts/:id | Update post | write |
| DELETE | /api/v1/posts/:id | Delete post | admin |

**Create/Update post body:**

```json
{
  "title": "Post Title",
  "slug": "post-slug",
  "content": "HTML content",
  "excerpt": "Brief summary",
  "cover": "https://example.com/cover.jpg",
  "status": "draft",
  "category_id": 1,
  "tag_names": ["tag1", "tag2"]
}
```

`status` options: `draft`, `published`, `archived`. If slug is omitted, it will be auto-generated from the title.

### Pages

| Method | Endpoint | Description | Required Permission |
|--------|----------|-------------|-------------------|
| GET | /api/v1/pages | List pages | read |
| GET | /api/v1/pages/:id | Get page detail | read |
| POST | /api/v1/pages | Create page | write |
| PUT | /api/v1/pages/:id | Update page | write |
| DELETE | /api/v1/pages/:id | Delete page | admin |

### Categories

| Method | Endpoint | Description | Required Permission |
|--------|----------|-------------|-------------------|
| GET | /api/v1/categories | List all categories | read |
| POST | /api/v1/categories | Create category | write |
| PUT | /api/v1/categories/:id | Update category | write |
| DELETE | /api/v1/categories/:id | Delete category | admin |

### Tags

| Method | Endpoint | Description | Required Permission |
|--------|----------|-------------|-------------------|
| GET | /api/v1/tags | List tags (paginated) | read |
| POST | /api/v1/tags | Create tag | write |
| PUT | /api/v1/tags/:id | Update tag | write |
| DELETE | /api/v1/tags/:id | Delete tag | admin |

### Media

| Method | Endpoint | Description | Required Permission |
|--------|----------|-------------|-------------------|
| GET | /api/v1/media | List media files | read |
| DELETE | /api/v1/media/:id | Delete media | admin |

### Comments

| Method | Endpoint | Description | Required Permission |
|--------|----------|-------------|-------------------|
| GET | /api/v1/comments | List comments (supports `status` filter) | read |
| DELETE | /api/v1/comments/:id | Delete comment | admin |

### Site Settings

| Method | Endpoint | Description | Required Permission |
|--------|----------|-------------|-------------------|
| GET | /api/v1/settings | Get all settings | read |
| PUT | /api/v1/settings | Update settings | admin |

## Operating Guidelines

1. **Read before write:** Always fetch existing data with GET before creating or updating. This prevents duplicates and helps the user make informed decisions.
2. **Confirm destructive actions:** Before DELETE operations, explicitly ask the user for confirmation. Show what will be deleted.
3. **Default to draft:** When creating posts, default `status` to `draft` unless the user explicitly asks to publish.
4. **Report results:** After each operation, summarize what was done (created/updated/deleted, IDs, titles).
5. **Handle errors gracefully:** If an API returns a non-zero code, explain the error to the user in plain language.
6. **Batch operations:** For bulk changes, make multiple API calls and report aggregate results.
7. **Respect permissions:** If a key lacks write/admin permission, inform the user that the operation requires elevated access.

## Full API Schema

For complete OpenAPI 3.0 specification, see: /api/ai/schema (relative to your CMS base URL)
