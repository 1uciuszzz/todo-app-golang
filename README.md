# todo-app-golang

go 语言实现的 todo restful api 应用

# 启动项目

```bash
go install
go run ./cmd/api/main.go
```

# Ubuntu 服务器上部署

## 1. 系统环境准备

首先，我们需要在 Ubuntu 服务器上安装必要的软件包。以下命令需要用 root 权限执行：

```bash
# 更新系统包
sudo apt update
sudo apt upgrade -y

# 安装必要的系统工具
sudo apt install -y curl git build-essential

# 安装 PostgreSQL
sudo apt install -y postgresql postgresql-contrib

# 安装 Go（使用最新的稳定版本）
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
rm go1.21.6.linux-amd64.tar.gz

# 将 Go 添加到系统路径
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
source ~/.bashrc
```

## 2. 配置 PostgreSQL

```bash
# 切换到 postgres 用户
sudo -i -u postgres

# 创建数据库和用户
psql -c "CREATE DATABASE todo_db;"
psql -c "CREATE USER todo_user WITH ENCRYPTED PASSWORD 'your_password';"
psql -c "GRANT ALL PRIVILEGES ON DATABASE todo_db TO todo_user;"

# 退出 postgres 用户
exit
```

## 3. 项目部署

### 3.1 创建部署目录和用户

```bash
# 创建服务用户
sudo useradd -m -s /bin/bash todo_api

# 创建应用目录
sudo mkdir -p /opt/todo-api
sudo chown todo_api:todo_api /opt/todo-api
```

### 3.2 部署项目代码

```bash
# 切换到 todo_api 用户
sudo su - todo_api

# 克隆项目代码
cd /opt/todo-api
git clone https://github.com/yourusername/todo-api.git .

# 安装依赖
go mod download
go mod verify
```

### 3.3 创建环境配置文件

```bash
# 创建并编辑 .env 文件
cat > /opt/todo-api/.env << EOF
SERVER_ADDRESS=:8080
DATABASE_URL=postgres://todo_user:your_password@localhost:5432/todo_db
EOF
```

## 4. 使用 Systemd 管理服务

创建系统服务配置文件：

```bash
sudo tee /etc/systemd/system/todo-api.service << EOF
[Unit]
Description=Todo API Service
After=network.target postgresql.service

[Service]
Type=simple
User=todo_api
WorkingDirectory=/opt/todo-api
ExecStart=/usr/local/go/bin/go run cmd/api/main.go
Restart=always
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF
```

启动并启用服务：

```bash
# 重新加载 systemd 配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start todo-api

# 设置开机自启
sudo systemctl enable todo-api

# 检查服务状态
sudo systemctl status todo-api
```

## 5. 配置 Nginx 反向代理（可选）

如果你想通过域名访问 API，可以安装和配置 Nginx：

```bash
# 安装 Nginx
sudo apt install -y nginx

# 创建 Nginx 配置文件
sudo tee /etc/nginx/sites-available/todo-api << EOF
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
    }
}
EOF

# 启用站点配置
sudo ln -s /etc/nginx/sites-available/todo-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

## 6. 防火墙配置

如果你启用了 UFW 防火墙，需要开放必要的端口：

```bash
# 开放 SSH 端口（如果尚未开放）
sudo ufw allow ssh

# 如果直接访问 API
sudo ufw allow 8080/tcp

# 如果使用 Nginx
sudo ufw allow 'Nginx Full'

# 启用防火墙
sudo ufw enable
```

## 7. 日志管理

查看服务日志：

```bash
# 查看实时日志
sudo journalctl -u todo-api -f

# 查看最近的日志
sudo journalctl -u todo-api -n 100
```

## 8. 备份策略

创建数据库备份脚本：

```bash
sudo tee /opt/todo-api/backup.sh << EOF
#!/bin/bash
BACKUP_DIR="/opt/todo-api/backups"
TIMESTAMP=\$(date +%Y%m%d_%H%M%S)
mkdir -p \$BACKUP_DIR

# 数据库备份
pg_dump -U todo_user todo_db > \$BACKUP_DIR/todo_db_\$TIMESTAMP.sql

# 保留最近 7 天的备份
find \$BACKUP_DIR -type f -mtime +7 -delete
EOF

sudo chmod +x /opt/todo-api/backup.sh
```

添加定时任务：

```bash
# 编辑 crontab
sudo crontab -e

# 添加每日备份任务（每天凌晨 2 点执行）
0 2 * * * /opt/todo-api/backup.sh
```
