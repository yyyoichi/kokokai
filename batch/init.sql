create table kyoki_day (kyoki_day_id integer primary key, date date);
create table kyoki(
    kyoki_id integer primary key,
    kyoki_day_id integer,
    freq integer
);
create table kyoki_item (
    kyoki_item_id integer primary key,
    kyoki_id integer,
    kyoki_day integer,
    word_id integer
);
create table word (word_id integer, word varchar(100));