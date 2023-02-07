-- Delete all previous roles
DELETE FROM roles;

-- Create multiple roles
INSERT INTO roles VALUES
  (DEFAULT, 'Admin', now(), now()),
  (DEFAULT, 'User', now(), now());