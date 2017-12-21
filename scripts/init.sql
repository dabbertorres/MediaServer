create database if not exists radio;

use radio;

create table if not exists users
(
    name     varchar(32) primary key,
    password binary(60) not null
);

create table if not exists sessions
(
    id        binary(16) primary key,
    user      varchar(32)   not null unique,
    ipAddr    varbinary(16) not null,
    userAgent varchar(32)   not null,
    expires   timestamp     not null,
    tunedTo   varchar(32),
    foreign key (user) references users (name)
        on delete set null
        on update cascade,
    foreign key (tunedTo) references stations (name)
        on delete set null
        on update cascade
);

create table if not exists songs
(
    title  varchar(32)  not null,
    artist varchar(32)  not null,
    source varchar(128) not null primary key
);

# load songs database, ignoring duplicates
load data infile 'songs.csv' ignore into table songs
fields terminated by ','
enclosed by '"'
lines terminated by '\n';

create table if not exists stations
(
    name     varchar(32) not null primary key,
    owner    binary(16)  not null unique,
    playlist json        not null
);
