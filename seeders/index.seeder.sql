-- Remove all previous users
DELETE FROM users;

-- Create multiple users
INSERT INTO users VALUES
  (DEFAULT, 'Admin', 'admin@gmail.com', '$2a$14$IlIhsldNRhUN/5SXXC1NeO.AU0YOlLQmeu7bg1y7tD2cornrAZ5ty', now(), now()),
  (DEFAULT, 'Huda', 'huda@gmail.com', '$2a$14$IlIhsldNRhUN/5SXXC1NeO.AU0YOlLQmeu7bg1y7tD2cornrAZ5ty', now(), now());

-- Delete all previous roles
DELETE FROM roles;

-- Create multiple roles
INSERT INTO roles VALUES
  (DEFAULT, 'Admin', now(), now()),
  (DEFAULT, 'User', now(), now());

-- Delete all previous permissions
DELETE FROM permissions;

-- Create multiple permissions
INSERT INTO permissions VALUES
  (DEFAULT, 'TODO', now(), now()),
  (DEFAULT, 'USER_MANAGEMENT', now(), now()),
  (DEFAULT, 'ROLE_MANAGEMENT', now(), now());