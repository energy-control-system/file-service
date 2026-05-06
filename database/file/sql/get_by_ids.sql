select id, file_name, file_size, bucket, url
from storage_files
where id = any ($1)
order by id
limit $2 offset $3;
