package input

// 输入系统初始化
func InputSysInit() {
	nowInputCli()
}

// api
func InputUpdate() { nowInputCli().input() }

func HasAnyKeyPressed() bool         { return nowInputCli().hasAnyKeyPressed }  // 检查是否有按键被按下
func HasAnyKeyReleased() bool        { return nowInputCli().hasAnyKeyReleased } // 检查是否有按键被松开
func PauseInput()                    { nowInputCli().pauseInput() }             // 暂停input接受，但不会清空bind
func ContinueInput()                 { nowInputCli().continueInput() }          // 继续input
func (k *actionKeyInfo) Check() bool { return nowInputCli().checkKey(k) }       // 单纯检查key，一般用于热键绑定
func ClearActionBuf()                { nowInputCli().clearActionBuf() }         // 清空所有按键缓存

// 绑定
func CheckAllMouseReleasedAction() bool { return nowInputCli().checkAllMouseReleasedAction() } // 检查鼠标左右键是否松开行为
func CheckAllMousePerssedAction() bool  { return nowInputCli().checkAllMousePerssedAction() }  // 检查鼠标左右键是否有按下行为

// 新建按键行为
func NewAction(keyType keyPressedType) *actionKeyInfo {
	return &actionKeyInfo{
		keyType: keyType,
		keys:    make([]inputKeyI, 0),
	}
}
