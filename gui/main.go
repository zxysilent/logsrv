package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

type AppWin struct {
	*walk.MainWindow
	teMsg         *walk.TextEdit
	btnTest       *walk.PushButton
	btnLsvm       *walk.PushButton
	btnLsct       *walk.PushButton
	btnClone      *walk.PushButton
	btnRetry      *walk.PushButton
	rbAuto        *walk.RadioButton
	rbClone       *walk.RadioButton
	title         string //标题
	newSuffix     string //新的后缀名
	oldSuffix     string //待删除后缀
	vmsPath       string //保存路径
	retryPath     string //重试路径
	concurrentNum int    //并发
}

var app *AppWin
var retrys []string   //重试列表
var wg sync.WaitGroup // 并发备份
var queue chan bool   //并发备份
var errRetry = errors.New("存在失败记录")

func main() {
	app = &AppWin{}
	app.title = "虚拟机备份管理"
	walk.App().SetProductName(app.title)
	walk.App().SetOrganizationName("zxysilent")
	icon, _ := walk.NewIconFromResourceId(3)
	err := MainWindow{
		Visible:  false,
		AssignTo: &app.MainWindow,
		Title:    app.title,
		Size:     Size{Width: 520, Height: 400},
		Font:     Font{Family: "微软雅黑", PointSize: 10},
		Icon:     icon,
		MenuItems: []MenuItem{
			Action{
				Text:        "使用说明",
				OnTriggered: app.OnTips,
			},
			Action{
				Text:        "关于",
				OnTriggered: app.OnAbout,
			},
		},
		Layout: VBox{},
		Children: []Widget{
			GroupBox{
				Title:  "循环周期",
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "后缀名", EllipsisMode: EllipsisEnd, MaxSize: Size{Width: 260, Height: 30}},
					Composite{
						Layout:  HBox{},
						MaxSize: Size{Width: 260, Height: 30},
						Children: []Widget{
							RadioButtonGroup{
								DataMember: "rb_set",
								Buttons: []RadioButton{
									{
										Name:     "rb_auto",
										Text:     "周期AUTO",
										Value:    "_autoc_",
										AssignTo: &app.rbAuto,
									},
									{
										Name:     "rb_clone",
										Text:     "周期CLONE",
										Value:    "_aclone_",
										AssignTo: &app.rbClone,
									},
								},
							},
						},
					},
					// PushButton{
					// 	MaxSize: Size{Width: 100, Height: 30},
					// 	Text:    "周期检测",
					// },
				},
			},
			GroupBox{
				Title:  "功能区",
				Layout: HBox{Spacing: 10, Alignment: AlignHCenterVCenter},
				Children: []Widget{
					Composite{
						Layout:  HBox{},
						MaxSize: Size{Width: 320, Height: 30},
						Children: []Widget{
							PushButton{
								Text:      "测试",
								AssignTo:  &app.btnTest,
								MaxSize:   Size{Width: 60, Height: 30},
								OnClicked: app.OnTest,
							},
							PushButton{
								Text:      "列出",
								AssignTo:  &app.btnLsvm,
								MaxSize:   Size{Width: 60, Height: 30},
								OnClicked: app.OnLsvm,
							},
							PushButton{
								Text:      "计数",
								AssignTo:  &app.btnLsct,
								MaxSize:   Size{Width: 60, Height: 30},
								OnClicked: app.OnLsct,
							},
							PushButton{
								Text:      "克隆",
								AssignTo:  &app.btnClone,
								MaxSize:   Size{Width: 60, Height: 30},
								OnClicked: app.OnClone,
							},
							PushButton{
								Text:      "重试",
								AssignTo:  &app.btnRetry,
								MaxSize:   Size{Width: 60, Height: 30},
								OnClicked: app.OnRetry,
							},
						},
					},
				},
			},
			TextEdit{AssignTo: &app.teMsg, VScroll: true, ReadOnly: true, Font: Font{Family: "微软雅黑", PointSize: 8}},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					// PushButton{
					// 	MinSize: Size{160, 30},
					// 	Text:    "打开windows服务管理程序",
					// 	OnClicked: func() {
					// 	},
					// },
					HSpacer{},
					PushButton{
						MinSize: Size{Width: 100, Height: 30},
						Text:    "清空日志",
						OnClicked: func() {
							app.teMsg.SetText("")
						},
					},
					PushButton{
						MinSize: Size{Width: 100, Height: 30},
						Text:    "关闭窗口",
						OnClicked: func() {
							walk.App().Exit(0)
						},
					},
				},
			},
		},
		OnSizeChanged: func() {
			app.SetSize(walk.Size(Size{Width: 520, Height: 400}))
		},
	}.Create()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	winLong := win.GetWindowLong(app.Handle(), win.GWL_STYLE)
	// 不能调整窗口大小，禁用最大化按钮
	win.SetWindowLong(app.Handle(), win.GWL_STYLE, winLong & ^win.WS_SIZEBOX & ^win.WS_MAXIMIZEBOX & ^win.WS_SIZEBOX)
	// 设置窗体生成在屏幕的正中间，并处理高分屏的情况
	// 窗体横坐标 = ( 屏幕宽度 - 窗体宽度 ) / 2
	// 窗体纵坐标 = ( 屏幕高度 - 窗体高度 ) / 2
	app.SetX((int(win.GetSystemMetrics(0)) - app.Width()) / 2 / app.DPI() * 96)
	app.SetY((int(win.GetSystemMetrics(1)) - app.Height()) / 2 / app.DPI() * 96)
	app.Show()
	app.Run()
}

func (app *AppWin) Log(args ...interface{}) {
	app.teMsg.AppendText(time.Now().Format("01/02 15:04:05] "))
	app.teMsg.AppendText(fmt.Sprint(args...))
	app.teMsg.AppendText("\r\n")
}
func (app *AppWin) OnTips() {
	app.teMsg.SetText("")
	app.Log("1.请选择管理平台已有的后缀,「警告」克隆或者重试会自动删除后缀虚拟机，\r\n2.点击功能区按钮\r\n* 测试-检测服务器连接情况\r\n* 列出-输出开机虚拟机信息\r\n* 计数-计算平台虚拟机数量\r\n* 克隆-根据配置文件进行克隆\r\n* 重试-对克隆失败的进行重新克隆\r\n")
}
func (app *AppWin) OnAbout() {
	// walk.MsgBox(app, "关于", "作者：曾祥银\n日期：20210322\n版本：1.0.0", walk.MsgBoxIconInformation)
	app.teMsg.SetText("")
	app.Log("作者：曾祥银\r\n日期：20210322\r\n版本：1.0.0")
}

// 测试链接
func (app *AppWin) OnTest() {

}

// 列出虚拟机
func (app *AppWin) OnLsvm() {

}
func (app *AppWin) OnLsct() {

}

func (app *AppWin) OnClone() {

}
func (app *AppWin) OnRetry() {

}

// isExist 文件存在
func (app *AppWin) isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// inVms 判断虚拟机是否在内
func (app *AppWin) inVms(dist string, arr []string) bool {
	for idx := range arr {
		if arr[idx] == dist {
			return true
		}
	}
	return false
}
