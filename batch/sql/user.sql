CREATE TABLE usr (
    pk serial primary key,
    id varchar(24),
    name varchar(20),
    email varchar(50),
    pass varchar(24),
    login_at timestamp,
    update_at timestamp,
    create_at timestamp DEFAULT CURRENT_TIMESTAMP
);