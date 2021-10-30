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
	action     CHAR(1)      NOT NULL,
	stamp      TIMESTAMP    NOT NULL DEFAULT NOW(),
	id         BIGINT       NOT NULL,
	name       VARCHAR(255) NOT NULL,
	slug       VARCHAR(255) NOT NULL,
	created_at TIMESTAMP    NOT NULL,
	updated_at TIMESTAMP    NOT NULL
);

CREATE TABLE entries_history
(
	action       CHAR(1)   NOT NULL,
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
	updated_at   TIMESTAMP NOT NULL
);

-- migrate:down

DROP TABLE IF EXISTS entries_history;
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS logs_history;
DROP TABLE IF EXISTS logs;
