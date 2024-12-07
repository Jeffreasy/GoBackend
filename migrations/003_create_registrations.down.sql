-- +migrate Down
DROP TABLE IF EXISTS registrations;
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS distance;
DROP TYPE IF EXISTS support_type; 