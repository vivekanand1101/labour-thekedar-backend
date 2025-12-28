# ADR 002: PostgreSQL Database

## Status
Accepted

## Context
We need a database to store:
- User accounts and authentication data
- Projects and their metadata
- Labour records
- Attendance/work day records
- Payment transactions

The system requires ACID compliance, relational data modeling, and good performance.

## Decision
We will use **PostgreSQL 15** as the primary database.

## Rationale

1. **ACID Compliance**: Full transaction support for financial data (payments)
2. **Relational Model**: Natural fit for our entity relationships (users -> projects -> labours)
3. **UUID Support**: Native UUID type for primary keys
4. **JSON Support**: JSONB for flexible metadata if needed
5. **Mature & Stable**: Battle-tested in production environments
6. **Open Source**: No licensing costs
7. **Excellent Go Support**: pgx driver is well-maintained and performant

## Schema Design Decisions

### Primary Keys
- Using UUIDs instead of auto-increment integers
- Provides better distribution and security (non-guessable IDs)

### Timestamps
- All tables include `created_at` and `updated_at`
- Using `TIMESTAMPTZ` for timezone awareness

### Enums
- Using PostgreSQL ENUM types for `work_status` and `payment_type`
- Provides type safety at database level

## Alternatives Considered

### MySQL
- Pros: Widely used, familiar
- Cons: Weaker JSON support, fewer advanced features

### MongoDB
- Pros: Flexible schema, good for prototyping
- Cons: No ACID by default, not ideal for financial data

### SQLite
- Pros: Simple, no server needed
- Cons: Not suitable for production concurrent access

## Consequences

### Positive
- Reliable storage for financial transactions
- Strong data integrity with foreign keys
- Excellent query performance with proper indexing
- Easy backup and replication options

### Negative
- Requires running a separate database server
- Schema migrations need careful planning
- Slightly more complex setup than SQLite
