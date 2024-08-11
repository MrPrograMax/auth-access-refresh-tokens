CREATE TABLE "user"
(
    id             uuid        not null unique primary key,
    email          varchar(64) not null unique,
    password       varchar(128) not null,
    refresh_token  varchar(128),
    ip varchar(15) not null,
    expires_at     timestamp
);