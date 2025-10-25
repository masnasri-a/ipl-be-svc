# API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Interactive Documentation
Swagger UI is available at: `http://localhost:8080/swagger/index.html`

## Authentication
Currently, the API does not require authentication. This can be added later as needed.

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

## Endpoints

### Health Check

#### GET /health
Check if the service is running.

**Response:**
```json
{
  "status": "ok",
  "message": "Server is running",
  "service": "IPL Backend Service"
}
```

### Menus

#### GET /menus/user/:user_id
Get list of menus accessible by a specific user ID.

**Path Parameters:**
- `user_id`: User ID (integer)

**SQL Query Used:**
```sql
SELECT DISTINCT ON (mm.document_id) mm.*
FROM up_users_role_lnk uurl
INNER JOIN role_menus_role_lnk rmrl ON rmrl.role_id = uurl.role_id
INNER JOIN role_menus_master_menu_lnk rmmml ON rmrl.role_menu_id = rmmml.role_menu_id
INNER JOIN master_menus mm ON rmmml.master_menu_id = mm.id
WHERE uurl.user_id = ?
ORDER BY mm.document_id, mm.id;
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Menus retrieved successfully",
  "data": [
    {
      "id": 1,
      "document_id": "mo5qqs8ezbruui07t91p6da8",
      "nama_menu": "Master Data",
      "kode_menu": "master-data",
      "urutan_menu": 1,
      "is_active": true,
      "published_at": null
    },
    {
      "id": 3,
      "document_id": "xhfm87n1dsc6a8u6bp4l7yj1",
      "nama_menu": "User Management",
      "kode_menu": "user-management",
      "urutan_menu": 2,
      "is_active": true,
      "published_at": "2025-10-23T15:16:28.206Z"
    }
  ]
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Invalid user ID",
  "error": "strconv.ParseUint: parsing \"invalid\": invalid syntax"
}
```

**Response (200 OK - No menus):**
```json
{
  "success": true,
  "message": "No menus found for this user",
  "data": []
}
```

## Error Codes

- `400 Bad Request`: Invalid request parameters
- `404 Not Found`: Resource not found
- `405 Method Not Allowed`: HTTP method not allowed
- `500 Internal Server Error`: Server error

## Examples with cURL

### Health Check
```bash
curl -X GET http://localhost:8080/api/v1/health
```

### Get menus for user ID 1
```bash
curl -X GET http://localhost:8080/api/v1/menus/user/1
```

### Get menus for user ID with error handling
```bash
curl -X GET http://localhost:8080/api/v1/menus/user/invalid
```

## Menu Response Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `id` | integer | Unique menu identifier | `1` |
| `document_id` | string | Menu document identifier | `"mo5qqs8ezbruui07t91p6da8"` |
| `nama_menu` | string | Display name of the menu | `"Master Data"` |
| `kode_menu` | string | Menu code/slug | `"master-data"` |
| `urutan_menu` | integer/null | Menu display order | `1` |
| `is_active` | boolean/null | Whether the menu is active | `true` |
| `published_at` | string/null | Publication timestamp | `"2025-10-23T15:16:28.206Z"` |

## Database Schema Expected

The API expects the following tables to exist in the PostgreSQL database:

### master_menus
```sql
CREATE TABLE master_menus (
    id SERIAL PRIMARY KEY,
    document_id VARCHAR(255),
    name VARCHAR(255),
    url VARCHAR(255),
    icon VARCHAR(255),
    order_num INTEGER,
    parent_id INTEGER,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

### Related tables for user-role-menu relationships
- `up_users_role_lnk` - Links users to roles
- `role_menus_role_lnk` - Links roles to role_menus
- `role_menus_master_menu_lnk` - Links role_menus to master_menus

## Business Logic

1. **User Validation**: User ID must be a valid positive integer
2. **Active Filtering**: Only active menus (`is_active = true`) are returned
3. **Distinct Results**: Uses `DISTINCT ON (document_id)` to avoid duplicate menu documents
4. **Ordered Results**: Results are ordered by `document_id` and `id` for consistent output

## Testing the API

1. **Start the server:**
   ```bash
   go run cmd/server/main.go
   ```

2. **Test health endpoint:**
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

3. **Test menu endpoint:**
   ```bash
   curl http://localhost:8080/api/v1/menus/user/1
   ```

4. **View Swagger documentation:**
   Open `http://localhost:8080/swagger/index.html` in your browser

## Rate Limiting

Currently, no rate limiting is implemented. This can be added as middleware if needed.

## Monitoring

The application logs all requests with:
- Request method and path
- Response status code
- Request duration
- Client IP
- User agent
- Request/response body size