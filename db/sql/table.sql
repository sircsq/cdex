CREATE TABLE collections(
    id SERIAL not null,
    name varchar(32) not null ,
    address varchar(65) not null,
    creator varchar(65) not null,
    chain smallint not null,
    visible smallint not null,
    status smallint not null,
    created_at timestamp not null,
    type smallint not null,
    tax smallint not null,
    symbol varchar(16) not null,
    currency varchar(16) not null,
    image varchar(128),
    background varchar(128),
    banner varchar(128),
    properties text,
    introduction varchar(256),
    description varchar(256),
    twitter varchar(256),
    instagram varchar(256),
    discord varchar(256),
    web varchar(256),
    PRIMARY KEY (id)
);

CREATE TABLE items(
    name varchar(32) not null,
    collection integer not null,
    token_id integer,
    creator varchar(65) not null,
    created_at timestamp not null,
    chain smallint not null,
    image varchar(128) not null,
    description varchar(128),
    properties varchar(128),
    primary key (collection, token_id)
);
