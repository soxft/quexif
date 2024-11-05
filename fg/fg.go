package fg

import (
	"flag"
)

var Mode string
var Path string
var DateTime string
var DateTpl = "2006-01-02 15.04.05"

// Parse 解析命令行参数
func Parse() {
	flag.StringVar(&Mode, "m", "qumagie", "操作模式: qumagie (QuMagie 备份照片处理), dir (指定文件夹批量修改 EXIF时间), dirDate (按照上级文件夹名称修改 EXIF 时间)")
	flag.StringVar(&Path, "p", "", "文件夹路径")
	flag.StringVar(&DateTime, "datetime", "", "日期时间")
	flag.StringVar(&DateTpl, "tpl", "2006-01-02 15.04.05", "日期时间模板, 默认为 '2006-01-02 15.04.05' 请参照 Golang 时间 layout 设置")

	flag.Parse()
}