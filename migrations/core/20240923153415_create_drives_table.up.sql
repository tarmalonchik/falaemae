drop table if exists drives;
create table drives
(
    id         uuid     default uuid_generate_v4() not null constraint drives_pk primary key,
    driver_id uuid not null,
    time timestamptz not null,
    direction text not null,
    car   int not null,
    price int not null,
    slots int not null,
    created_at timestamptz default now()
);

create unique index drives_id_uindex
    on drives (id);

drop table if exists cars;
create table cars
(
    id         serial constraint cars_pk primary key,
    name text not null

);
create unique index cars_name_uindex
    on cars (name);