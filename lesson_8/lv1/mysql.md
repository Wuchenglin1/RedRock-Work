```shell
docker pull mysql #获取镜像

#运行容器，需要做数据挂载
#安装启动mysql，需要配置密码 
#$ docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:tag

#启动容器
[root@VM-0-15-centos ceshi]# docker run -d -p 3310:3306 -v /home/mysql/conf:/etc/mysql/conf.d -v /home/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root --name mysql01 mysql
```

