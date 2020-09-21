package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"kingwrcy/dnspod-client/api"
	"kingwrcy/dnspod-client/model"
)

type DomainModel struct {
	walk.TableModelBase
	items []model.Domain
}

func (m *DomainModel) RowCount() int {
	return len(m.items)
}
func (m *DomainModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Name
	case 1:
		return item.Records
	}
	panic("unexpected col")
}

type RecordModel struct {
	walk.TableModelBase
	items []model.Record
}

func (m *RecordModel) RowCount() int {
	return len(m.items)
}
func (m *RecordModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.ID

	case 1:
		return item.Type

	case 2:
		return item.Name

	case 3:
		return item.Value
	}
	panic("unexpected col")
}

func loadRecord() {
	index := domainBox.CurrentIndex()
	if index <= 0 {
		index = 0
	}
	var records = api.GetRecordList(domainModel.items[index].ID)
	recordModel.items = records
	recordModel.PublishRowsReset()
}

var (
	mainWin   *walk.MainWindow
	domainBox *walk.TableView
	recordBox *walk.TableView

	domainModel *DomainModel
	recordModel *RecordModel

	record *model.Record
	dlg    *walk.Dialog
)

func add() {
	record = new(model.Record)
	showDialog()
}

func edit() {
	index := recordBox.CurrentIndex()
	if index < 0 {
		walk.MsgBox(dlg, "错误", "请先选一条记录", walk.MsgBoxIconError)
		return
	}
	record = &recordModel.items[index]
	showDialog()
}

func removeRecord() {
	recordIndex := recordBox.CurrentIndex()
	domainIndex := domainBox.CurrentIndex()
	if recordIndex < 0 {
		walk.MsgBox(mainWin, "错误", "请先选一条记录", walk.MsgBoxIconError)
		return
	}
	if domainIndex < 0 {
		walk.MsgBox(mainWin, "错误", "请先选一个域名", walk.MsgBoxIconError)
		return
	}
	ok := api.RemoveRecord(domainModel.items[domainIndex].ID, recordModel.items[recordIndex].ID)
	if ok {
		loadRecord()
	}
}

func saveOrUpdateRecord() {
	if domainBox.CurrentIndex() < 0 {
		walk.MsgBox(mainWin, "错误", "请先选一个域名", walk.MsgBoxIconError)
		return
	}
	domainId := domainModel.items[domainBox.CurrentIndex()].ID

	ok := false
	if record.ID == 0 {
		ok = api.SaveRecord(domainId, *record)
	} else {
		ok = api.ModifyRecord(domainId, *record)
	}
	if ok {
		loadRecord()
		dlg.Accept()
	}
}

func showDialog() {

	var acceptPB *walk.PushButton
	var cancelPB *walk.PushButton
	var dialog = Dialog{
		Icon: "3",
	}
	var db *walk.DataBinder

	if record.ID == 0 {
		dialog.Title = "新增记录"
	} else {
		dialog.Title = "编辑记录"
	}
	dialog.DataBinder = DataBinder{
		AssignTo:       &db,
		Name:           "record",
		DataSource:     record,
		ErrorPresenter: ToolTipErrorPresenter{},
	}
	dialog.MinSize = Size{Width: 300, Height: 200}
	dialog.Layout = VBox{}
	dialog.DefaultButton = &acceptPB
	dialog.CancelButton = &cancelPB
	dialog.AssignTo = &dlg

	childrens := []Widget{
		Composite{
			Layout: Grid{Columns: 2},
			Children: []Widget{
				Label{
					Text: "记录类型:",
				},
				ComboBox{
					Value: Bind("Type"),
					Model: []string{"A", "CNAME", "MX", "TXT", "NS", "AAAA", "SRV", "CAA"},
				},
				Label{
					Text: "名称:",
				},
				LineEdit{
					Text:      Bind("Name"),
				},
				Label{
					Text: "内容:",
				},
				TextEdit{
					MinSize: Size{Height: 50},
					MaxSize: Size{Height: 250},
					Text: Bind("Value"),
				},
			},
		},
		Composite{
			Layout: HBox{},
			Children: []Widget{
				HSpacer{},
				PushButton{
					AssignTo: &acceptPB,
					Text:     "保存",
					OnClicked: func() {
						if err := db.Submit(); err != nil {
							return
						}
						saveOrUpdateRecord()

					},
				},
				PushButton{
					AssignTo:  &cancelPB,
					Text:      "取消",
					OnClicked: func() { dlg.Cancel() },
				},
			},
		},
	}
	dialog.Children = childrens
	dialog.Run(mainWin)
}

func main() {

	token := api.GetLoginToken()
	if token == "" {
		return
	}

	recordModel = &RecordModel{
		items: []model.Record{},
	}

	domainModel = &DomainModel{
		items: []model.Domain{},
	}

	err := MainWindow{
		Icon:     "3",
		AssignTo: &mainWin,
		Title:    "Dnspod",
		Size:     Size{Width: 900, Height: 600},
		Layout: Grid{
			Columns: 2, Spacing: 10,
		},
		Children: []Widget{
			Composite{
				Layout:  VBox{},
				MinSize: Size{Width: 160},
				Children: []Widget{
					TextLabel{
						Text: "域名列表:",
					},
					TableView{
						AssignTo:       &domainBox,
						MultiSelection: false,
						Columns: []TableViewColumn{
							{Title: "域名"},
							{Title: "记录数量", Width: 70},
						},
						Model:                 domainModel,
						OnCurrentIndexChanged: loadRecord,
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					ToolBar{
						Orientation: Horizontal,
						ButtonStyle: ToolBarButtonImageBeforeText,
						Items: []MenuItem{
							Action{
								Image:       "5",
								Text:        "添加",
								OnTriggered: add,
							}, Action{
								Image:       "9",
								Text:        "编辑",
								OnTriggered: edit,
							}, Action{
								Image:       "7",
								Text:        "删除",
								OnTriggered: removeRecord,
							}, Action{
								Image:       "11",
								Text:        "刷新",
								OnTriggered: loadRecord,
							},
						},
					},
					TableView{
						MinSize:          Size{Width: 600},
						AssignTo:         &recordBox,
						AlternatingRowBG: true,
						ColumnsOrderable: true,
						MultiSelection:   false,
						Columns: []TableViewColumn{
							{Title: "ID"},
							{Title: "记录类型"},
							{Title: "记录"},
							{Title: "内容", Width: 320},
						},
						Model: recordModel,
					},
				},
			},
		},
	}.Create()
	if err == nil {
		var domains = api.GetDomainList()
		domainModel.items = domains
		_ = domainBox.SetModel(domainModel)

		if len(domains) > 0 {
			var records = api.GetRecordList(domains[0].ID)
			recordModel.items = records
		}
		_ = domainBox.SetCurrentIndex(0)
		mainWin.Run()
	} else {
		fmt.Printf("error:%s\n", err.Error())
	}
}
