# MySQL Example for tqla

This example demonstrates how to use the `tqla` library with MySQL database.

## Set up the environment

### 1. Pull and run MySQL Docker container

```shell
# Pull the latest MySQL image
docker pull mysql:latest

# Start MySQL container
docker run --name mysql_test -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=mysql-test-db -p 3306:3306 -d mysql:latest
```
