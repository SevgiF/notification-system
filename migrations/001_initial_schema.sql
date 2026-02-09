CREATE TABLE IF NOT EXISTS status (
    code INT PRIMARY KEY,
    description VARCHAR(50) NOT NULL
);

INSERT IGNORE INTO status (code, description) VALUES (0, 'deleted'), (1, 'created'), (2, 'canceled'), (3, 'sent'), (4, 'failed'), (5, 'processing');

CREATE TABLE IF NOT EXISTS notification (
    id INT AUTO_INCREMENT PRIMARY KEY,
    recipient VARCHAR(255) NOT NULL,
    channel VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    priority VARCHAR(10) NOT NULL,
    status INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (status) REFERENCES status(code)
);

CREATE INDEX idx_notification_status_priority_created ON notification(status, priority, created_at);
