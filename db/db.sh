


docker exect -it postgresDB /bin/bash
createdb --username=root --owner=root simple_bank
psql simple_bank
\q
dropdb simple_bank

docker exect -it postgresDB createdb --username=root --owner=root simple_bank
docker exect -it postgresDB psql -U root simple_bank