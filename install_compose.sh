sudo curl -L "https://get.daocloud.io/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" \
-o /usr/bin/docker-compose

#centos7放在/usr/local/bin/docker-compose
#对二进制文件应用可执行权限：

sudo chmod +x /usr/bin/docker-compose