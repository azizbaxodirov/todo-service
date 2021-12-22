CREATE TABLE IF NOT EXISTS todos(
    id uuid Primary Key,
    assignee VARCHAR(50),
    title VARCHAR(50),
    summary VARCHAR(50),
    deadline  timestamp not null,
    status VARCHAR(50)
);
