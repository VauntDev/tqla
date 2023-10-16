### Set up the environment

```shell
# pulling the latest postgres image
docker pull postgres

# starting postgres container
docker run -e POSTGRES_PASSWORD=root -e POSTGRES_USER=root --name psql_test -p 5432:5432 -d postgres:latest
```