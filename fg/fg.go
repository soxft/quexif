package fg

import (
	"flag"
)

var Mode string
var Path string
var DateTime string
var DateTpl = "2006-01-02 15.04.05"
var Force bool
var SkipSafeQA bool

// Parse 解析命令行参数
func Parse() {
	flag.StringVar(&Mode, "m", "read", "操作模式: qumagie (QuMagie 备份照片处理), dir (指定文件夹批量修改 EXIF时间), dir_date (按照路径推导时间), read (读取目录或文件的 EXIF 时间信息)")
	flag.StringVar(&Path, "p", "", "文件夹路径")
	flag.StringVar(&DateTime, "d", "", "日期时间")
	flag.StringVar(&DateTpl, "t", "2006-01-02 15.04.05", "日期时间模板, 默认为 '2006-01-02 15.04.05' 请参照 Golang 时间 layout 设置, 不适用于 QuMagie 模式")
	flag.BoolVar(&Force, "f", false, "强制执行, 不会检查是否已经有日期")
	flag.BoolVar(&SkipSafeQA, "skip", false, "跳过安全询问, 直接执行")

	flag.Parse()
}
