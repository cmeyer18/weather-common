CREATE TABLE IF NOT EXISTS device(
    id          varchar(255) primary key,
    userId      varchar(255),
    apnsToken   varchar(255)
)
