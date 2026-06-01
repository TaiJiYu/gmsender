package gmsender

import (
	"gmsender/pkg/asset"
	gametime "gmsender/pkg/game_time"
	"gmsender/pkg/input"
	"gmsender/pkg/netfinder"
	statemachine "gmsender/pkg/state_machine"
	"gmsender/pkg/ui"
	"gmsender/utils"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type gmsenderCli struct {
	sizex, sizey                   int
	canvas                         *ui.CanvasUi
	screenX, screenY               int  // 窗口位置
	mousePressPosX, mousePressPosY int  // 鼠标按下的位置
	screenDroping                  bool // 是否正在被拖动

	closeButton *ui.ButtonUi // 关闭按钮
	smallButton *ui.ButtonUi // 缩小按钮
	reInButton  *ui.ButtonUi // 重启按钮
	topCanvas   *ui.CanvasUi // 顶部右侧栏位

	midCanvas       *ui.CanvasUi // 中部区域的画布，主要功能区
	choseFileButton *ui.ButtonUi // 选择文件按钮

	filesVbox *ui.VerticalBox // 文件列表

	idText *ui.TextUi // 本机id

	smallState statemachine.StateMachineWithDrawI // 缩小状态机

	smallCacheImg    *ebiten.Image
	smallWeekdayText *ui.TextUi // 缩小后的文本组件
	lastMouseTick    int64      // 上一次点击的帧数

	isClose bool
}

var (
	sendercli  *gmsenderCli
	senderOnce sync.Once
)

const (
	loadingState = iota // 加载状态
	bigState            // 放大状态
	toSmallState        // 正在缩小
	smallState          // 已经缩小了
	toBigState          // 正在放大

	smallTime    = 200 * time.Millisecond // 200毫秒用于缩放
	smallTimeSec = 0.2                    // 200毫秒用于缩放
)

func NewGMSender() *gmsenderCli {
	gametime.InitGameTimer()
	ui.Init()
	input.InputSysInit()
	input.ContinueInput()

	InitFileList()

	senderOnce.Do(func() {
		sendercli = &gmsenderCli{
			sizex:            utils.LogicalSizeX,
			sizey:            utils.LogicalSizeY,
			smallCacheImg:    ebiten.NewImage(utils.SmallLogicalSize, utils.SmallLogicalSize),
			smallWeekdayText: ui.NewStaticTextUi(utils.WeekStr(), ui.SmallSize, utils.NewPointT1(utils.SmallLogicalSize).Divf1(2), utils.MM, color.White),
			lastMouseTick:    -1,
		}
		sendercli.smallCacheImg.DrawImage(asset.SmallBallImg(), utils.CopyDrawImageOp)
		sendercli.smallWeekdayText.Draw(sendercli.smallCacheImg)

		sendercli.canvas = ui.NewCoreRectCanvasUi(utils.ZeroPoint, utils.LL, backColor, 0)
		sendercli.canvas.AddKid(ui.NewSizeBox(utils.NewPoint(utils.LogicalSizeX, utils.LogicalSizeY)))

		sendercli.topCanvas = ui.NewEmptyCanvasUi(utils.ZeroPoint, utils.LL, 20)
		// 顶部栏位分左右部分，左侧是客户端id，右侧是缩小和关闭
		hBox := ui.NewHorizontalBox(20)
		sendercli.idText = hBox.AddKid(ui.NewStaticTextUiAsKid("ID:???", ui.SmallSize, idTextColor)).(*ui.TextUi)
		sendercli.topCanvas.AddKid(hBox)

		coloseImgX := float64(asset.CloseImg().Bounds().Dx())

		sendercli.closeButton = ui.NewButton(asset.CloseImg(), utils.NewPoint(utils.LogicalSizeX, 0).Sub(utils.NewPoint(20, -20)).Sub(utils.NewPoint(asset.CloseImg().Bounds().Dx()/2, -asset.CloseImg().Bounds().Dy()/2)), utils.MM, gametime.BigTimerType)
		sendercli.closeButton.SetCheckKey(input.GameMainReleasedAction, func(bu *ui.ButtonUi) {
			sendercli.isClose = true
		})

		sendercli.smallButton = ui.NewButton(asset.SmallImg(), utils.NewPoint(utils.LogicalSizeX, 0).SubX(coloseImgX+10).Sub(utils.NewPoint(20, -20)).Sub(utils.NewPoint(asset.SmallImg().Bounds().Dx()/2, -asset.SmallImg().Bounds().Dy()/2)), utils.MM, gametime.BigTimerType)
		sendercli.smallButton.SetCheckKey(input.GameMainReleasedAction, func(bu *ui.ButtonUi) {
			sendercli.smallState.Go(toSmallState)
		})

		sendercli.reInButton = ui.NewButton(asset.ReinImg(), utils.NewPoint(utils.LogicalSizeX, 0).SubX((coloseImgX+10)*2).Sub(utils.NewPoint(20, -20)).Sub(utils.NewPoint(asset.ReinImg().Bounds().Dx()/2, -asset.ReinImg().Bounds().Dy()/2)), utils.MM, gametime.BigTimerType)
		sendercli.reInButton.SetCheckKey(input.GameMainReleasedAction, func(bu *ui.ButtonUi) {
			sendercli.smallState.Go(loadingState)
		})

		sendercli.midCanvas = ui.NewEmptyCanvasUi(utils.NewPoint(utils.LogicalSizeX/2, 60), utils.ML, 0)

		vBox := ui.NewVerticalBox(10)
		choseCanvas := vBox.AddKid(ui.NewRoundLerpRectCanvasUiAsKid(choiseColor, choiseColor, 10).LockSize(utils.NewPoint(fileListX, 20)))
		sendercli.choseFileButton = ui.NewButtonByCanvas(choseCanvas.(*ui.CanvasUi), choiseColor, choiseColor, gametime.BigTimerType) // 选择文件公开按钮
		sendercli.choseFileButton.SetCheckKey(input.GameMainReleasedAction, func(bu *ui.ButtonUi) {
			netfinder.PublicFile(utils.OpenWinChooseFile())
		})
		choseCanvas.(*ui.CanvasUi).AddKid(ui.NewStaticTextUiAsKid("十 选择文件公开给别人", ui.SmallSize, choiseTextColor).AddSpaceToSizeX(utils.LogicalSizeX - 40*2))
		filescanvas := ui.NewCoreRectCanvasUiAsKid(fileListBackColor, 10).LockSize(utils.NewPoint(fileListX, utils.LogicalSizeY-20-60-20-10-40)) // 文件列表框

		sendercli.filesVbox = ui.NewVerticalBox(10)

		filescanvas.AddKid(sendercli.filesVbox)
		vBox.AddKid(filescanvas)
		sendercli.midCanvas.AddKid(vBox)

		sendercli.filesVbox.SetSlider(filescanvas.Size().Y-10*2, backColor)

		// 缩小状态机
		sendercli.smallState = statemachine.NewStateMachineWithDraw(gametime.BigTimerType)

		loadingTexts := []*ui.TextUi{
			ui.NewStaticTextUi(".", ui.BigSize, utils.LogicalSize.Divf1(2), utils.MM, color.White),
			ui.NewStaticTextUi("..", ui.BigSize, utils.LogicalSize.Divf1(2), utils.MM, color.White),
			ui.NewStaticTextUi("...", ui.BigSize, utils.LogicalSize.Divf1(2), utils.MM, color.White),
			ui.NewStaticTextUi("网络错误，启动失败..", ui.MidSize, utils.LogicalSize.Divf1(2), utils.MM, color.White),
		}
		loadingCanvas := ui.NewCoreRectCanvasUi(utils.ZeroPoint, utils.LL, backColor, 0).LockSize(utils.LogicalSize)
		loadingChangeSec := (500 * time.Millisecond).Seconds()
		loadingTimeLimit := (30 * time.Second).Seconds()
		// loadingCheckTImeMin := (5 * time.Second).Seconds()
		isFaile := false
		// 加载状态
		sendercli.smallState.NewState(loadingState).SetEnterFunc(func() {
			isFaile = false
			netfinder.Init(sendercli.refreshFiles, sendercli.typeChange, func() {
				sendercli.smallState.Go(bigState)
			})
		}).SetExitFunc(func() {
			sendercli.idText.SetText("ID:" + netfinder.Id())
		}).BindUD(func() {
			gametime.BigTimeRun()
			sendercli.moveScreen()

			t := sendercli.smallState.ReadStateLastTimeSec()

			if t > loadingTimeLimit {
				// 启动超时
				input.InputUpdate()
				sendercli.closeButton.Update(utils.NewPoint(ebiten.CursorPosition()))
			}
			// if t > loadingCheckTImeMin {
			// 	if done, err := netfinder.IsInitDone(); err != nil {
			// 		isFaile = true
			// 	} else if done {
			// 		sendercli.smallState.Go(bigState)
			// 	}
			// }
		}, func(screen *ebiten.Image) {
			screen.Clear()
			loadingCanvas.Draw(screen)

			if isFaile {
				loadingTexts[3].Draw(screen)
			} else {
				t := int(sendercli.smallState.ReadStateLastTimeSec()/loadingChangeSec) % 3
				loadingTexts[t].Draw(screen)
			}

			if sendercli.smallState.ReadStateLastTimeSec() > loadingTimeLimit || isFaile {
				// 启动超时
				sendercli.closeButton.Draw(screen)
			}
		})

		sendercli.smallState.NewState(bigState).SetEnterFunc(func() {
			sendercli.sizex = utils.LogicalSizeX
			sendercli.sizey = utils.LogicalSizeY
			ebiten.SetWindowSize(utils.LogicalSizeX, utils.LogicalSizeY)
		}).BindUD(func() {
			gametime.BigTimeRun()
			input.InputUpdate()
			sendercli.moveScreen()

			if input.MouseWheelDownAction.Check() {
				// 滑动值增加
				sendercli.filesVbox.Slider(20)
			} else if input.MouseWheelUpAction.Check() {
				// 滑动值减少
				sendercli.filesVbox.Slider(-20)
			}
			checkPos := utils.NewPoint(ebiten.CursorPosition())

			sendercli.closeButton.Update(checkPos)
			sendercli.smallButton.Update(checkPos)
			sendercli.reInButton.Update(checkPos)
			sendercli.choseFileButton.Update(checkPos)

			if sendercli.filesVbox.CheckMouseInSlider(checkPos) {
				fileListCli.Update(checkPos)
			}
		}, func(screen *ebiten.Image) {
			screen.Clear()
			sendercli.canvas.Draw(screen)
			sendercli.topCanvas.Draw(screen)
			sendercli.closeButton.Draw(screen)
			sendercli.smallButton.Draw(screen)
			sendercli.reInButton.Draw(screen)

			sendercli.midCanvas.Draw(screen)
		})
		sendercli.smallState.NewState(toSmallState).BindUD(func() {
			gametime.BigTimeRun()

			t := sendercli.smallState.ReadStateLastTimeSec() / smallTimeSec
			sizeX := utils.Lerp(utils.LogicalSizeX, utils.SmallLogicalSize, t)
			sizeY := utils.Lerp(utils.LogicalSizeY, utils.SmallLogicalSize, t)
			ebiten.SetWindowSize(int(sizeX), int(sizeY))
		}, func(screen *ebiten.Image) {
			screen.Clear()
			sendercli.canvas.Draw(screen)
			sendercli.topCanvas.Draw(screen)
			sendercli.closeButton.Draw(screen)
			sendercli.smallButton.Draw(screen)
			sendercli.reInButton.Draw(screen)

			sendercli.midCanvas.Draw(screen)
		})

		sendercli.smallState.NewState(smallState).SetEnterFunc(func() {
			sendercli.sizex = utils.SmallLogicalSize
			sendercli.sizey = utils.SmallLogicalSize
			ebiten.SetWindowSize(utils.SmallLogicalSize, utils.SmallLogicalSize)
			ebiten.SetTPS(utils.SmallTPS)
		}).SetExitFunc(func() {
			ebiten.SetTPS(utils.BigTPS)
		}).BindUD(func() {
			gametime.SmallTimeRun()
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if sendercli.lastMouseTick < 0 {
					sendercli.lastMouseTick = ebiten.Tick()
				} else if now := ebiten.Tick(); now-sendercli.lastMouseTick <= utils.DoubleClickTime {
					// 在双击范围内
					sendercli.lastMouseTick = -1
					sendercli.smallState.Go(toBigState)
					sendercli.screenDroping = false
					return
				} else {
					sendercli.lastMouseTick = now
				}
			}
			sendercli.moveScreen()

		}, func(screen *ebiten.Image) {
			screen.DrawImage(sendercli.smallCacheImg, utils.CopyDrawImageOp)
		})
		sendercli.smallState.NewState(toBigState).BindUD(func() {
			gametime.BigTimeRun()
			t := sendercli.smallState.ReadStateLastTimeSec() / smallTimeSec
			sizeX := utils.Lerp(utils.SmallLogicalSize, utils.LogicalSizeX, t)
			sizeY := utils.Lerp(utils.SmallLogicalSize, utils.LogicalSizeY, t)
			ebiten.SetWindowSize(int(sizeX), int(sizeY))
		}, func(screen *ebiten.Image) {
			screen.Clear()
			sendercli.canvas.Draw(screen)
			sendercli.topCanvas.Draw(screen)
			sendercli.closeButton.Draw(screen)
			sendercli.smallButton.Draw(screen)
			sendercli.reInButton.Draw(screen)

			sendercli.midCanvas.Draw(screen)
		})

		sendercli.smallState.SToSWithTimeLimit(toSmallState, smallState, smallTime)
		sendercli.smallState.SToSWithTimeLimit(toBigState, bigState, smallTime)

		sendercli.smallState.Go(loadingState)
	})

	return sendercli
}

// 刷新文件展示
func (s *gmsenderCli) refreshFiles(files []netfinder.File) {
	fileListCli.refreshFiles(files)
}

// 客户端身份转换
func (s *gmsenderCli) typeChange(isMaster bool) {
	if isMaster {
		sendercli.idText.SetColor(idMasterColor)
	} else {
		sendercli.idText.SetColor(idNodeColor)
	}
}

// 添加一个文件到列表展示
func appendFileCmp(f *ui.CanvasUi) {
	sendercli.filesVbox.AddKid(f)
}

// 刷新公开文件列表展示
func delFileCmp(f *ui.CanvasUi) {
	sendercli.filesVbox.DelKid(f)
}

// 移动窗口检测
func (s *gmsenderCli) moveScreen() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// 按下
		s.screenX, s.screenY = ebiten.WindowPosition()
		s.mousePressPosX, s.mousePressPosY = ebiten.CursorPosition()
		scale := ebiten.Monitor().DeviceScaleFactor()
		s.mousePressPosX = int(float64(s.mousePressPosX)/scale) + s.screenX
		s.mousePressPosY = int(float64(s.mousePressPosY)/scale) + s.screenY
		s.screenDroping = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// 松开
		s.screenDroping = false
	}
	if s.screenDroping {
		screenX, screenY := ebiten.WindowPosition()
		mousePosx, mousePosy := ebiten.CursorPosition()
		scale := ebiten.Monitor().DeviceScaleFactor()
		x := int(float64(mousePosx)/scale) + screenX - s.mousePressPosX
		y := int(float64(mousePosy)/scale) + screenY - s.mousePressPosY
		ebiten.SetWindowPosition(x+s.screenX, y+s.screenY)
	}
}

func (s *gmsenderCli) Update() error {
	if s.isClose {
		return ebiten.Termination
	}
	s.smallState.Update()

	return nil
}

func (s *gmsenderCli) Draw(screen *ebiten.Image) {
	s.smallState.Draw(screen)
}

func (s *gmsenderCli) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return s.sizex, s.sizey
	// return utils.LogicalSizeX, utils.LogicalSizeY
}
