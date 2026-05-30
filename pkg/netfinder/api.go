package netfinder

func Init() {
	defaultFinder()
}

// 本机id
func Id() string {
	return id
}

// 下载文件
func DownLoadFile()

// 关闭网络
func Close() {
	defaultFinder().closeNetFinder()
}
