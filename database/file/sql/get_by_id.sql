select id, file_name, file_size, url
from storage_files
where id = $1;
