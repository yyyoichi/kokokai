create table kyokiday (
    pk serial primary key,
    date date,
    create_at timestamp DEFAULT CURRENT_TIMESTAMP
);
create table kyoki(
    pk serial primary key,
    kyokiday integer,
    freq integer
);
create table kyokiitem (
    pk serial primary key,
    kyokiday integer,
    kyoki integer,
    word integer
);
create table word (code serial primary key, word varchar(100));