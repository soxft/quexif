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


## Thinks

- [go-exif](//github.com/dsoprea/go-exif/v3)