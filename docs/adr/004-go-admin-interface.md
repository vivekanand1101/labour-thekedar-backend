# ADR 004: go-admin for Admin Interface

## Status
Accepted

## Context
We need an admin interface for:
- Managing users (view, block/unblock)
- Viewing all projects and labours
- Payment reconciliation
- System monitoring

The admin interface should:
- Require minimal frontend development
- Integrate with existing Go backend
- Provide CRUD operations
- Support role-based access

## Decision
We will use **go-admin** (GoAdminGroup/go-admin) for the admin interface.

## Rationale

1. **Native Go Integration**: Works directly with Gin framework
2. **Auto CRUD Generation**: Automatically generates forms and tables
3. **Built-in RBAC**: Role-based access control included
4. **Theme Support**: Multiple UI themes available
5. **Plugin System**: Extensible for custom functionality
6. **Active Development**: Regularly maintained and updated

## Features We'll Use

### Data Tables
- Sortable, filterable tables for all entities
- Pagination built-in
- Export to CSV/Excel

### Forms
- Auto-generated create/edit forms
- Validation support
- File upload if needed

### Dashboard
- Custom dashboard with statistics
- Total projects, labours, payments
- Recent activity

### Access Control
- Super admin for full access
- Read-only roles for viewing data

## Configuration
```go
adminPlugin := admin.NewAdmin(dataList)
adminPlugin.AddGenerator("users", GetUserTable)
adminPlugin.AddGenerator("projects", GetProjectTable)
adminPlugin.AddGenerator("labours", GetLabourTable)
adminPlugin.AddGenerator("payments", GetPaymentTable)
```

## Alternatives Considered

### Custom React Admin
- Pros: Full control, modern UI
- Cons: Significant development time, separate codebase

### AdminLTE Template
- Pros: Popular, good looking
- Cons: Manual integration, no CRUD generation

### Gin-admin
- Pros: Lighter weight
- Cons: Less features, smaller community

## Consequences

### Positive
- Rapid admin interface development
- No frontend expertise required
- Consistent look and feel
- Built-in security features

### Negative
- Less flexibility than custom solution
- Learning curve for go-admin specific concepts
- UI may not match exact design requirements
- Additional dependency in the project

## Port Configuration
- Main API: Port 8080
- Admin Interface: Port 9033
