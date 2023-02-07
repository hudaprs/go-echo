-- Delete all previous permissions
DELETE FROM permissions;

-- Create multiple permissions
INSERT INTO permissions VALUES
  (DEFAULT, 'TODO', now(), now()),
  (DEFAULT, 'USER_MANAGEMENT', now(), now()),
  (DEFAULT, 'ROLE_MANAGEMENT', now(), now());