insert into storage_files (file_name, file_size, bucket, url)
values ($1, $2, $3, $4)
returning id, file_name, file_size, url;
