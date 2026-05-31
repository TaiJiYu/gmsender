package netfinder

var (
	filesCallback      func(files []File)
	typeChangeCallback func(isMaster bool)
)

// filesCallback文件同步回调，公开文件有更新时调用
func Init(filesCallbackFunc func(files []File), typeChangeCallbackFunc func(isMaster bool)) {
	filesCallback = filesCallbackFunc
	typeChangeCallback = typeChangeCallbackFunc
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
	defaultFinder().publicFile(filename)
}

// 删除公开文件
func DelPublicFile(file File) {
	defaultFinder().delPublicFile(file)
}

// 下载文件
func DownLoadFile() {}

// 关闭网络
func Close() {
	defaultFinder().closeNetFinder()
}
