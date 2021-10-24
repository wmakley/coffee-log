-- migrate:up
INSERT INTO logs (name, slug, created_at, updated_at)
VALUES ('Default', 'default', NOW(), NOW());

-- migrate:down
DELETE FROM logs WHERE slug = 'default';
