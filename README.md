# Quexif

> 自动为 Qumagie 备份的手机照片, 添加 Exif 信息, 以便于在相册中按照时间线查看


## 原理

在 Qumagie 备份时, 需要取消选择 原始文件名称, 这样备份的文件名会采用日期格式, 例如 `2019-01-01 12.00.00.jpg`

此工具将读取时间信息, 并将其写入 Exif 的 DateTime 和 DateTimeOriginal 信息，此时 Qumagie 将会将此时间作为 照片的拍摄日期。

## 使用

```shell
    go run main.go /path/to/your/photos
```