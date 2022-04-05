create database todos;

create table todos (
    id INT,
    assignee VARCHAR(128),
    title VARCHAR(64),
    summary VARCHAR(64),
    deadline TIMESTAMPTZ default CURRENT_TIMESTAMP,
    status VARCHAR(64)
);