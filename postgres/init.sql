CREATE DATABASE supertodo_users;

\c supertodo_users;

DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username VARCHAR(255) UNIQUE
);

CREATE DATABASE supertodo_todos;

\c supertodo_todos;

DROP TABLE IF EXISTS todos;

CREATE TABLE todos (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title VARCHAR(255),
    todo_date TIMESTAMP,
    body VARCHAR
);

CREATE DATABASE supertodo_combined;

\c supertodo_combined;

DROP TABLE IF EXISTS combined;

CREATE TABLE combined (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INT,
    todo_id INT
    -- CONSTRAINT fk_user
    --     FOREIGN KEY(user_id) REFERENCES users(id),
    -- CONSTRAINT fk_todo
    --     FOREIGN KEY(todo_id) REFERENCES todos(id)
);