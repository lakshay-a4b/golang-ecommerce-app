-- Down migration: Drops users table and related indexes
DROP TABLE IF EXISTS users CASCADE;