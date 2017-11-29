create database if not exists radio;

use radio;

create table if not exists users
(
    name varchar(64),
    password binary(60) not null,
    primary key (name)
);

create table if not exists sessions
(
    id binary(16),
    user varchar(64),
    tunedTo varchar(64),
    primary key (id),
    foreign key (user) references users (name)
        on delete set null
        on update cascade,
    foreign key (tunedTo) references stations (name)
        on delete set null
        on update cascade
);

create table if not exists songs
(
    name varchar(64) not null,
    artist varchar(64) not null,
    genre varchar(32),
    source varchar(128) not null,
    primary key (source)
);

create table if not exists stations
(
    name varchar(64) not null,
    owner binary(16) not null,
    playlistFile varchar(16) not null
);
