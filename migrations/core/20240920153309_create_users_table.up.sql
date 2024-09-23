create extension if not exists "uuid-ossp";
drop table if exists users;
create table users
(
    id         uuid     default uuid_generate_v4() not null constraint users_pk primary key,
    telegram_user     text,
    chat_id           integer                   not null,
    created_at        timestamptz default now() not null
);

create unique index users_chat_id_uindex
    on users (chat_id);

create unique index users_id_uindex
    on users (id);

create unique index users_telegram_user_uindex
    on users (telegram_user);