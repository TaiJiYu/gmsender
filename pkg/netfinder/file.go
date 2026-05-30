package netfinder

// 文件结构
type file struct {
	Ip       string `json:"ip"`        // 文件所属ip
	Port     string `json:"port"`      // 用于下载该文件的对应端口
	Id       string `json:"id"`        // 文件所属id，用于判读是否是自己的
	FileName string `json:"file_name"` // 文件名
}
