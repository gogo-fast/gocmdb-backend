# gocmdb-backend

## Description

frontend code go to `->` [gocmdb-front](https://github.com/gogo-fast/gocmdb-front)



## Quick start with docker

- Docker version

  ```shell
  $ docker -v
  Docker version 19.03.1, build 74b1e89
  ```

- Clone code

  ```shell
  git clone https://github.com/gogo-fast/gocmdb-backend.git
  ```

- Build images

  ```shell
  cd gocmdb-backend
  docker build --no-cache -t cmdb-bin:v0.1 -f dockerfiles/gocmdb .
  docker build --no-cache -t cmdb-db:v0.1 -f dockerfiles/mysql dockerContexts/mysql/
  ```

- Run

  start mysql (`cmdb-db`), change `192.168.10.100` to `your_host_ip`

  ````shell
  docker run --name cmdb_db -d -p 192.168.10.100:3306:3306 cmdb-db:v0.1
  ````

  make sure `cmdb_db` start success

  ```shell
  docker ps -a | grep cmdb_db
  ```

  start `cmdb-bin`

  ```shell
  docker run --name cmdb_bg -d -p 192.168.10.100:8000:8000 cmdb-bin:v0.1
  ```

- Forgot password

  Generate a new password with `genpass`

  ```shell
  $ docker exec -it cmdb_bg /cmdb/bin/genpass -p 123456  # 123456 is new pass
  123456:75a0fc76c396f590c5e63c9cc58d260e
  ```

  Update password

  ```shell
  docker exec -it cmdb_db /bin/bash -c '''mysql -ucmdb -p123456 -e "use gocmdb; update users set password = \"123456:75a0fc76c396f590c5e63c9cc58d260e\" where username = \"admin\" "; '''  
  ```

  

