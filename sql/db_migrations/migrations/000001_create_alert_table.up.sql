CREATE TABLE IF NOT EXISTS alerts(
    id       varchar(255) primary key,
    payload  jsonb not null
)