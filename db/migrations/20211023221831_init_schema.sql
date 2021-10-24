-- migrate:up

CREATE TABLE logs
(
	id         BIGSERIAL PRIMARY KEY,
	name       VARCHAR(255)        NOT NULL,
	slug       VARCHAR(255) UNIQUE NOT NULL,
	created_at TIMESTAMP           NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP           NOT NULL DEFAULT NOW()
);

CREATE TABLE entries
(
	id           BIGSERIAL PRIMARY KEY,
	log_id       BIGINT       NOT NULL REFERENCES logs (id),
	coffee       VARCHAR(255) NOT NULL,
	water        VARCHAR(255),
	method       VARCHAR(255),
	grind        VARCHAR(255),
	tasting      VARCHAR(4000),
	addl_notes   VARCHAR(4000),
	coffee_grams INT,
	water_grams  INT,
	created_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
	updated_at   TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_entries_log_id ON entries (log_id);

CREATE TABLE logs_history
(
	history_id BIGSERIAL PRIMARY KEY,
	action     CHAR         NOT NULL,
	stamp      TIMESTAMP    NOT NULL DEFAULT NOW(),
	id         BIGINT       NOT NULL,
	name       VARCHAR(255) NOT NULL,
	slug       VARCHAR(255) NOT NULL,
	created_at TIMESTAMP    NOT NULL,
	updated_at TIMESTAMP    NOT NULL,
);

CREATE INDEX idx_logs_history_id ON logs_history (id);

CREATE TABLE entries_history
(
	history_id   BIGSERIAL PRIMARY KEY,
	action       CHAR      NOT NULL,
	stamp        TIMESTAMP NOT NULL DEFAULT NOW(),
	id           BIGINT    NOT NULL,
	log_id       BIGINT    NOT NULL,
	coffee       VARCHAR(255),
	water        VARCHAR(255),
	method       VARCHAR(255),
	grind        VARCHAR(255),
	tasting      VARCHAR(4000),
	addl_notes   VARCHAR(4000),
	coffee_grams INT,
	water_grams  INT,
	created_at   TIMESTAMP NOT NULL,
	updated_at   TIMESTAMP NOT NULL,
);

CREATE INDEX idx_entries_history_id ON entries_history (id);

-- CREATE FUNCTION process_logs_audit() RETURNS TRIGGER AS $logs_audit$
-- 	BEGIN
--         IF (TG_OP = 'DELETE') THEN
-- 			INSERT INTO logs_history SELECT 'D', now(), OLD.*;
-- 		ELSIF (TG_OP = 'UPDATE') THEN
-- 			INSERT INTO logs_history SELECT 'U', now(), NEW.*;
-- 		ELSIF (TG_OP = 'INSERT') THEN
-- 			INSERT INTO logs_history SELECT 'I', now(), NEW.*;
-- 		END IF;
-- 		RETURN NULL;
-- 	END;
-- $logs_audit$ LANGUAGE plpgsql;
--
-- CREATE TRIGGER logs_audit
-- 	AFTER INSERT OR UPDATE OR DELETE ON logs
-- 	FOR EACH ROW EXECUTE FUNCTION process_logs_audit();
--
-- CREATE FUNCTION process_entries_audit() RETURNS TRIGGER AS $entries_audit$
-- 	BEGIN
--         IF (TG_OP = 'DELETE') THEN
-- 			INSERT INTO entries_history SELECT 'D', now(), OLD.*;
-- 		ELSIF (TG_OP = 'UPDATE') THEN
-- 			INSERT INTO entries_history SELECT 'U', now(), NEW.*;
-- 		ELSIF (TG_OP = 'INSERT') THEN
-- 			INSERT INTO entries_history SELECT 'I', now(), NEW.*;
-- 		END IF;
-- 		RETURN NULL;
-- 	END
-- $entries_audit$ LANGUAGE plpgsql;
--
-- CREATE TRIGGER entries_audit
-- 	AFTER INSERT OR UPDATE OR DELETE ON entries
-- 	FOR EACH ROW EXECUTE FUNCTION process_entries_audit();

-- migrate:down

DROP TRIGGER entries_audit;
DROP TRIGGER logs_audit;
DROP FUNCTION process_entries_audit;
DROP FUNCTION process_logs_audit;
DROP TABLE entries_history;
DROP TABLE entries;
DROP TABLE logs_history;
DROP TABLE logs;
