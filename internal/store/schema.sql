CREATE TABLE IF NOT EXISTS boards (
 id VARCHAR(80) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
 title VARCHAR(120) NOT NULL,
 token_hash CHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
 revision BIGINT NOT NULL DEFAULT 0,
 state JSON NOT NULL,
 created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS events (
 board_id VARCHAR(80) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
 revision BIGINT NOT NULL,
 operation_id VARCHAR(80) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
 fingerprint CHAR(64) NOT NULL,
 payload JSON NOT NULL,
 PRIMARY KEY(board_id, revision),
 UNIQUE KEY operation_once(board_id, operation_id),
 FOREIGN KEY(board_id) REFERENCES boards(id)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS snapshots (
 board_id VARCHAR(80) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
 revision BIGINT NOT NULL,
 state JSON NOT NULL,
 PRIMARY KEY(board_id,revision),
 FOREIGN KEY(board_id) REFERENCES boards(id)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS outbox (
 id BIGINT AUTO_INCREMENT PRIMARY KEY,
 board_id VARCHAR(80) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
 revision BIGINT NOT NULL,
 payload JSON NOT NULL,
 published_at TIMESTAMP(6) NULL,
 UNIQUE KEY event_once(board_id,revision),
 KEY pending(published_at,id),
 FOREIGN KEY(board_id,revision) REFERENCES events(board_id,revision)
) ENGINE=InnoDB;
CREATE TABLE IF NOT EXISTS audit_events (
 board_id VARCHAR(80) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
 revision BIGINT NOT NULL,
 payload JSON NOT NULL,
 consumed_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 PRIMARY KEY(board_id,revision),
 FOREIGN KEY(board_id,revision) REFERENCES events(board_id,revision)
) ENGINE=InnoDB;
