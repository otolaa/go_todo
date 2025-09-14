# ⚙️ setting for todo_bot

rename `.env.example` on `.env` and update settings for you bot

```
TOKEN=7969376429:CBEWj2Bm5-if5G194gT7r2kAQjn3MORyhVE
URL_BOT=https://ysqwp-83-219-149-215.a.free.pinggy.link
DB_DSN=host=localhost user=www dbname=go_todo_bot password=314159 port=5432 sslmode=disable
DEBUG=false
```

## 🌵 webhook in telegram bot

```
go mod init
go mod tidy
```

## 🍐 webhook start pinggy.io

🍎 add url webhook in .env

```
ssh -p 443 -R0:127.0.0.1:8080 -L4300:localhost:4300 free.pinggy.io
```

## 🍏 Run

```
killall -9 go
go run .
```

## 🌶️ nginx

```
go mod init go_todo
go build

systemctl restart nginx
```

## 🍎 systemd

```
sudo nano /etc/systemd/system/go_todo.service
```

add this command
```
[Unit]
Description=go_todo

[Service]
User=www
Group=www
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory=/home/www/goproject/go_todo/
ExecStart=/home/www/goproject/go_todo/go_todo

[Install]
WantedBy=multi-user.target
```

command for start
```
sudo systemctl start go_todo
sudo systemctl enable go_todo
sudo systemctl status go_todo
sudo systemctl restart go_todo

sudo systemctl stop go_todo
sudo systemctl disable go_todo
```

## 🍎 PostgreSQL

add table
```
sudo -u postgres psql
CREATE DATABASE go_todo;
ALTER DATABASE go_todo OWNER TO www;
```

Important: if no tables are found, you need to make sure that the user is in the correct database — to switch, you can use the command

\c database_name

command in psql for add database_name

```
sudo -u postgres psql
\c go_todo
\dt table_name_users
SELECT * FROM table_name_users;
```