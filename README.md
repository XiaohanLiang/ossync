# ossync

## 使用场景

- 主机端:
  - 你的服务运行在**很多主机**上,这个服务需要很多依赖
  - 依赖被放在一个统一的文件夹下,通过`docker -v`的方式挂载进服务
  - 这些依赖可能会**增加/更新**一些,你需要在主机上同步
  
- OSS端:
  - 你使用OSS来保存这些依赖
  - 每个依赖包都被打包压缩成了`tar.gz`的格式

## 安装与使用

### 安装它

- <u>直接下载二进制包使用</u>: [点这里](https://github.com/XiaohanLiang/ossync/releases)
- 通过源码编译:  
  - 需要 `go version > 1.11` 并且 `GO111MODULE=on`
  - 进入项目根目录执行 `go isntall main.go`

### 使用它

ossync被做成一个命令行工具的形式,包含六个参数
- `-e`  OSS的EndPoint
- `-r`  OSS的地域
- `-a`  你的AccessKey
- `-s`  你的SecretKey
- `-b`  Bucket名称
- `-p`  本地用于存放依赖的路径

``` bash
ossync -e s3.cn-north-1.jdcloud-oss.com -r cn-north-1 -a ABC***AccessKey -s ABC***SecretKey -b OSS -p /Users/mac/ossync_test/
```

### 工作原理

1. 简单来说,我们通过调用SDK获取OSS上所有文件的meta信息-最后修改时间, 生成一个map = map[file_name] = last_modified_time
2. 前往本地用于同步的路径, 读取所有文件, 以及他们的最后修改时间,也生成这样一个文件名以及最后修改时间的map
3. 拿着OSS的map对照本地map: 所有远程有,但是本地没有的,都会被认为需要下载
4. 本地远程都有但是最后修改时间对不上的, 说明需要更新,也要下载
5. 下载需要的包,解压, **并且修改新文件的最后修改时间**

### Q.A

- 为什么不用rsync做同步
  - 不喜欢主从模式
- 为什么不用ETag或者MD5作为依据?
  - 本地依赖=文件夹 & 远程依赖=tar.gz
- 我不是京东云可以吗? 
  - 但凡支持AWS-SDK的OSS都可以哦
  

<div align=center><img width="200" height="200" src="forkandstar.jpeg"/></div>
