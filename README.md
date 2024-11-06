# QuExif

> 自动为 QuMagie 备份的手机照片, 添加 Exif 信息, 以便于在相册中按照时间线查看


## 原理

在 QuMagie 备份时, 需要取消选择 原始文件名称, 这样备份的文件名会采用日期格式, 例如 `2019-01-01 12.00.00.jpg`

此工具将读取时间信息, 并将其写入 Exif 的 DateTime 和 DateTimeOriginal 信息，此时 QuMagie 将会将此时间作为 照片的拍摄日期。

## 使用

### 直接在本地运行, 以 Windows 为例

1. 通过 SMB 等方式将图片目录挂载到本地电脑 假设为 Z:\
2. 在 Release 目录下载对应的 二进制文件, 通常为 `quexif-windows-amd64-{{version}}.exe`
3. 在命令行中执行

```shell
$ quexif-windows-amd64-{{version}}.exe -p Z:\
```

### 在 Qnap 中运行

> 假设 你的备份照片路径在 `/share/Public/Photo`

1. 通过 SSH 登录 Qnap
2. 下载对应的 二进制文件, 通常为 `quexif-linux-amd64-{{version}}`
3. 将二进制文件上传到 Qnap 的 `/share/Public` 目录
4. 在 SSH 中执行

```shell
$ sudo -s 
    
$ chmod +x /share/Public/quexif-linux-amd64-{{version}}

$ /share/Public/quexif-linux-amd64-{{version}} -p /share/Public/Photo
```

## 其他支持项

> 您可以使用 ./quexif -h 查看所有支持的参数
```shell
Usage of quexif:
  -d string
        日期时间
  -f    强制执行, 不会检查是否已经有日期
  -m string
        操作模式: qumagie (QuMagie 备份照片处理), dir (指定文件夹批量修改 EXIF时间), dir_date (按照路径推导时间) (default "qumagie")
  -p string
        文件夹路径
  -skip
        跳过安全询问, 直接执行
  -t string
        日期时间模板, 默认为 '2006-01-02 15.04.05' 请参照 Golang 时间 layout 设置, 不适用于 QuMagie 模式 (default "2006-01-02 15.04.05")

```

- 批量修改某个目录及其子目录下的所有图片为指定时间

```shell
$ ./quexif -m dir -d '2024-11-23' -t '2006-01-02' -p ./pics

# -m dir 表示修改目录下的所有图片
# -d '2024-11-23' 表示修改为 2024-11-23
# -t '2006-01-02' 表示时间格式为 2006-01-02
# -p ./pics 表示目录为 ./pics
```

- 批量修改某个目录及其子目录下的所有图片, 按照设定的时间模板尝试推导时间

```shell
# 例如您的目录结构为
.
├── 2022-11-23
│   ├── IMG_0001.JPG
│   ├── IMG_0002.JPG
│   ├── IMG_0003.JPG
├── 2023-11-23
│   ├── IMG_0004.JPG
├── 2024-11-23
│   ├── IMG_0005.JPG
│   ├── IMG_0006.JPG

此时您可以执行如下脚本, 脚本将自动解析文件夹名称, 并将其作为时间写入 Exif

$ ./quexif -m dir_date -t '2006-01-02' -p ./pics

# -m dir_date 表示修改目录下的所有图片, 并按照文件名推导时间
# -t '2006-01-02' 表示时间格式为 2006-01-02
# -p ./pics 表示目录为 ./pics
```

## Thinks

- [go-exif](//github.com/dsoprea/go-exif/v3)