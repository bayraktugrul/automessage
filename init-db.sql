DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS recipients;

CREATE TABLE recipients (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(15) NOT NULL UNIQUE,
    name VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL CHECK (LENGTH(content) <= 160),
    recipient_id INTEGER NOT NULL REFERENCES recipients(id),
    is_sent BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP,
    message_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_messages_is_sent ON messages(is_sent);
CREATE INDEX idx_messages_created_at ON messages(created_at);

INSERT INTO recipients (phone_number, name, created_at, updated_at)
VALUES 
    ('+905551234567', 'Recipient 1', NOW(), NOW()),
    ('+905551234568', 'Recipient 2', NOW(), NOW()),
    ('+905551234569', 'Recipient 3', NOW(), NOW()),
    ('+905551234570', 'Recipient 4', NOW(), NOW());

INSERT INTO messages (content, recipient_id, is_sent, sent_at, created_at, updated_at)
VALUES 
    ('test message 1', 1, false, '2025-03-01 21:55:52', '2025-03-01 21:55:52', '2025-03-01 21:55:52'),
    ('test message 2', 1, false, NULL, '2025-03-01 21:55:52', '2025-03-01 21:55:52'),
    ('test message 3', 3, true, '2025-03-01 21:55:52', '2025-03-01 21:55:52', '2025-03-01 21:55:52'),
    ('test message 4', 4, false, '2025-03-01 21:55:52', '2025-03-01 21:55:52', '2025-03-01 21:55:52');
