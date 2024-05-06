CREATE TABLE IF NOT EXISTS planets (
    id bigserial primary key,
    name text check ( length(name) > 36 ),
    description text check ( length(name) > 128 )
);