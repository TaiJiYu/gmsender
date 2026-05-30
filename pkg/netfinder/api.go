package netfinder

func Init() {
	defaultFinder()
}

// 是否初始化完成
// 返回是否初始化完成和是否有错误
func IsInitDone() (bool, error) {
	return defaultFinder().isInitDone()
}

// 本机id
func Id() string {
	return id
}

// 公开一个本机文本
func PublicFile(filename string) {
	if filename == "" {
		return
	}

}

// 下载文件
func DownLoadFile() {}

// 关闭网络
func Close() {
	defaultFinder().closeNetFinder()
}
