-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;
DROP TRIGGER IF EXISTS update_labours_updated_at ON labours;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS otp_codes;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS work_days;
DROP TABLE IF EXISTS project_labours;
DROP TABLE IF EXISTS labours;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;

-- Drop enum types
DROP TYPE IF EXISTS payment_type;
DROP TYPE IF EXISTS work_status;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
