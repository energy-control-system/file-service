-- +goose Up
create table if not exists storage_files
(
    id        int primary key generated always as identity,
    file_name text not null,
    file_size int  not null,
    bucket    text not null,
    url       text not null -- Ссылка на скачивание
);

-- +goose Down
drop table if exists storage_files;
