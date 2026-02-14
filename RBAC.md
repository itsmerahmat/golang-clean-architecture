# RBAC (Role-Based Access Control) Feature

## Overview

This implementation adds comprehensive Role-Based Access Control (RBAC) to the Golang Clean Architecture application. RBAC allows you to manage user permissions through roles and control access to different resources.

## Features

- **Role Management**: Create, read, update, and delete roles
- **Permission Management**: Create, read, update, and delete permissions
- **User Role Assignment**: Assign and remove roles from users
- **Access Control**: Middleware for role and permission-based authorization
- **Default Roles**: Pre-configured admin and user roles
- **Default Permissions**: Pre-configured permissions for all resources

## Database Schema

The RBAC implementation adds the following tables:

### Tables

1. **roles**: Stores role information
   - `id`: Unique identifier
   - `name`: Role name (unique)
   - `description`: Role description
   - `created_at`, `updated_at`: Timestamps

2. **permissions**: Stores permission information
   - `id`: Unique identifier
   - `name`: Permission name (unique, e.g., "user:read")
   - `resource`: Resource name (e.g., "user", "contact")
   - `action`: Action name (e.g., "read", "write", "delete")
   - `description`: Permission description
   - `created_at`, `updated_at`: Timestamps

3. **user_roles**: Junction table for user-role relationships
   - `user_id`: Foreign key to users table
   - `role_id`: Foreign key to roles table

4. **role_permissions**: Junction table for role-permission relationships
   - `role_id`: Foreign key to roles table
   - `permission_id`: Foreign key to permissions table

## Default Configuration

### Default Roles

1. **admin**: Administrator with full access to all resources
2. **user**: Regular user with limited access

### Default Permissions

The system includes permissions for the following resources:
- **user**: read, write, delete
- **contact**: read, write, delete
- **address**: read, write, delete
- **role**: read, write, delete

### Permission Mapping

- **admin role**: Has all permissions
- **user role**: Has read/write/delete permissions for contacts and addresses, and read permission for users

## API Endpoints

### Role Management (Admin Only)

```
GET    /api/roles                      - List all roles
GET    /api/roles/:roleId              - Get role by ID
POST   /api/roles                      - Create new role
PUT    /api/roles/:roleId              - Update role
DELETE /api/roles/:roleId              - Delete role
```

### Permission Management (Admin Only)

```
GET    /api/permissions                - List all permissions
GET    /api/permissions/:permissionId  - Get permission by ID
POST   /api/permissions                - Create new permission
PUT    /api/permissions/:permissionId  - Update permission
DELETE /api/permissions/:permissionId  - Delete permission
```

### User Role Management (Admin Only)

```
GET    /api/users/:userId/roles             - Get user's roles
POST   /api/users/:userId/roles/:roleId     - Assign role to user
DELETE /api/users/:userId/roles/:roleId     - Remove role from user
```

## Usage Examples

### Creating a Role

```bash
curl -X POST http://localhost:3000/api/roles \
  -H "Authorization: <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "role-moderator",
    "name": "moderator",
    "description": "Moderator with limited admin access",
    "permissions": ["perm-user-read", "perm-contact-read", "perm-contact-write"]
  }'
```

### Assigning a Role to a User

```bash
curl -X POST http://localhost:3000/api/users/john@example.com/roles/role-admin \
  -H "Authorization: <admin-token>"
```

### Creating a Permission

```bash
curl -X POST http://localhost:3000/api/permissions \
  -H "Authorization: <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "perm-report-read",
    "name": "report:read",
    "resource": "report",
    "action": "read",
    "description": "Read reports"
  }'
```

## Middleware Usage

### RequireRole

Checks if the user has any of the specified roles:

```go
import "golang-clean-architecture/internal/delivery/http/middleware"

// Require admin role
adminOnly := middleware.RequireRole(logger, "admin")
app.Post("/api/admin/action", adminOnly, controller.AdminAction)

// Require admin or moderator role
adminOrMod := middleware.RequireRole(logger, "admin", "moderator")
app.Post("/api/moderate", adminOrMod, controller.ModerateAction)
```

### RequirePermission

Checks if the user has any of the specified permissions:

```go
// Require user:write permission
canWriteUser := middleware.RequirePermission(logger, "user:write")
app.Post("/api/users/create", canWriteUser, controller.CreateUser)

// Require contact:read or contact:write permission
canAccessContact := middleware.RequirePermission(logger, "contact:read", "contact:write")
app.Get("/api/contacts", canAccessContact, controller.GetContacts)
```

### RequireAllRoles

Checks if the user has all specified roles:

```go
// Require both admin and auditor roles
requireBoth := middleware.RequireAllRoles(logger, "admin", "auditor")
app.Post("/api/audit", requireBoth, controller.AuditAction)
```

### RequireAllPermissions

Checks if the user has all specified permissions:

```go
// Require all permissions
requireAll := middleware.RequireAllPermissions(logger, "user:write", "user:delete")
app.Delete("/api/users/bulk", requireAll, controller.BulkDeleteUsers)
```

## Migration

To apply the RBAC database changes, run:

```bash
migrate -database "mysql://root:@tcp(localhost:3306)/golang_clean_architecture?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations up
```

## Code Structure

The RBAC implementation follows the clean architecture pattern:

```
internal/
├── entity/
│   ├── role_entity.go          # Role entity
│   ├── permission_entity.go    # Permission entity
│   └── user_entity.go          # Updated with roles relationship
├── model/
│   ├── auth.go                 # Updated with roles/permissions
│   ├── role_model.go           # Role request/response models
│   ├── permission_model.go     # Permission request/response models
│   └── converter/
│       ├── role_converter.go
│       └── permission_converter.go
├── repository/
│   ├── role_repository.go
│   ├── permission_repository.go
│   └── user_repository.go      # Updated with role methods
├── usecase/
│   ├── role_usecase.go
│   ├── permission_usecase.go
│   └── user_usecase.go         # Updated to load roles
└── delivery/
    └── http/
        ├── middleware/
        │   └── rbac_middleware.go
        ├── role_controller.go
        ├── permission_controller.go
        └── user_role_controller.go
```

## Security Considerations

1. **Admin Protection**: All role and permission management endpoints are protected by admin role requirement
2. **Token-Based Authentication**: All RBAC endpoints require valid authentication token
3. **Cascading Deletes**: Deleting a role or permission automatically removes associated relationships
4. **Permission Validation**: Permissions are validated before being assigned to roles

## Best Practices

1. **Principle of Least Privilege**: Assign users only the roles they need
2. **Role Hierarchy**: Create roles with varying levels of access (e.g., admin, moderator, user)
3. **Permission Granularity**: Use fine-grained permissions for better control
4. **Regular Audits**: Regularly review user roles and permissions
5. **Default Role**: Consider assigning a default "user" role to new registrations

## Future Enhancements

Potential improvements for the RBAC system:

1. Role hierarchy with inheritance
2. Time-based role assignments
3. Resource-level permissions (e.g., user can only edit their own contacts)
4. Permission groups
5. Audit logging for role/permission changes
6. Role templates
7. Dynamic permission evaluation
