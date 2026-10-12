DROP TABLE IF EXISTS tasks;

CREATE TABLE tasks(
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       completed BOOLEAN DEFAULT FALSE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tasks (title, description, completed)
VALUES
    ('task1', 'desc1', FALSE),
    ('task2', 'desc2', FALSE),
    ('task3', 'desc3', FALSE);