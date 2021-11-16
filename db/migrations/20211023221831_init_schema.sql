-- migrate:up

CREATE TABLE users
(
	id           BIGSERIAL PRIMARY KEY,
	display_name VARCHAR(255)        NOT NULL,
	username     VARCHAR(100) UNIQUE NOT NULL,
	password     VARCHAR(255)        NOT NULL,
	time_zone    VARCHAR(100),
	created_at   TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at   TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE logs
(
	id         BIGSERIAL PRIMARY KEY,
	user_id    BIGINT              NOT NULL,
	slug       VARCHAR(255) UNIQUE NOT NULL,
	created_at TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_logs_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX index_logs_on_user_id ON logs (user_id);

CREATE TABLE log_entries
(
	id            BIGSERIAL PRIMARY KEY,
	log_id        BIGINT       NOT NULL REFERENCES logs (id),
	entry_date    TIMESTAMP    NOT NULL,
	coffee        VARCHAR(255) NOT NULL,
	water         VARCHAR(255),
	coffee_grams  INTEGER,
	water_grams   INTEGER,
	brew_method   VARCHAR(255),
	grind_notes   VARCHAR(255),
	tasting_notes VARCHAR(4000),
	addl_notes    VARCHAR(4000),
	deleted_at    TIMESTAMP,
	created_at    TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at    TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_log_entries_log_id FOREIGN KEY (log_id) REFERENCES logs (id)
);

CREATE INDEX index_log_entries_on_log_id_and_entry_date
	ON log_entries (log_id, entry_date)
	WHERE deleted_at IS NULL;

-- migrate:down

DROP TABLE IF EXISTS log_entries;
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS users;
