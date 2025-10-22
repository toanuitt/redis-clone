# 🚀 Hướng Dẫn Cài Đặt Redis Trên Ubuntu

## 🧩 1. Cài đặt gói cần thiết và thêm kho Redis chính thức

```bash
sudo apt-get install lsb-release curl gpg
curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg
sudo chmod 644 /usr/share/keyrings/redis-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list
sudo apt-get update
sudo apt-get install redis
```
2. sudo systemctl stop redis-server.service
3. redis-cli (with port 6379)