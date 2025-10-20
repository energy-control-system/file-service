select id, file_name, file_size, bucket, url
from storage_files
where id = $1;
