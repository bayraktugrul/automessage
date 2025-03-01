CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL CHECK (LENGTH(content) <= 160),
    recipient_phone VARCHAR(15) NOT NULL,
    is_sent BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_messages_is_sent ON messages(is_sent);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

INSERT INTO messages (content, recipient_phone) VALUES
('test message 1', '+905551234567'),
('test message 2', '+905551234568'),
('test message 3', '+905551234569'),
('test message 4', '+905551234570');
