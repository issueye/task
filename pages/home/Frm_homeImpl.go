package home

import (
	"fmt"
	"strconv"
	"strings"
	"task/internal/bdb"
	"task/internal/global"
	"task/internal/task"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/win"
)

func ShowHomeForm(owner vcl.IComponent) {
	if Frm_job == nil {
		Frm_job = NewFrm_job(owner)
	}

	Frm_job.list = make([]*bdb.Task, 0)
	Frm_job.Show()

	Frm_job.SetPosition(types.PoMainFormCenter)
	Frm_job.SetFormStyle(types.FsStayOnTop)

	Frm_job.GetData()

	width := vcl.Screen.Width()
	height := vcl.Screen.Height()
	Frm_job.SetLeft((width - Frm_job.Width()) / 2)
	Frm_job.SetTop((height - Frm_job.Height()) / 2)

	win.SetWindowPos(Frm_job.Handle(), win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)

}

// ::private::
type TFrm_jobFields struct {
	list []*bdb.Task
}

func (f *TFrm_job) OnFormCreate(sender vcl.IObject) {
	f.SetEvents()
}

func (f *TFrm_job) OnFormShow(sender vcl.IObject) {
}

func (f *TFrm_job) SetEvents() {
	f.TForm.SetOnShow(f.OnFormShow)
	f.Btn_add.SetOnClick(f.OnBtn_addClick)
	f.Btn_query.SetOnClick(f.OnBtn_queryClick)
	f.Table_data.SetOnButtonClick(f.OnTable_dataButtonClick)
}

func (f *TFrm_job) OnBtn_queryClick(sender vcl.IObject) {
	f.GetData()
}

func (f *TFrm_job) OnBtn_addClick(sender vcl.IObject) {
	ShowJobForm(f, nil)
	f.GetData()
}

func (f *TFrm_job) OnTable_dataButtonClick(sender vcl.IObject, aCol, aRow int32) {
	// 第五列 状态
	if aCol == 5 {
		// 更新状态
		if f.Table_data.Cells(5, aRow) == "1" {
			f.Table_data.SetCells(5, aRow, "0")
		} else {
			f.Table_data.SetCells(5, aRow, "1")
		}

		// 更新数据
		statusStr := f.Table_data.Cells(5, aRow)
		status, _ := strconv.Atoi(statusStr)
		f.list[aRow-1].Status = status
		bdb.GetBdb().UpdateTask(f.list[aRow-1])

		var err error
		if status == 1 {
			err = task.GetTaskCron().AddFunc(f.list[aRow-1], func() {
				global.Logger.Sugar().Debugf("任务执行: %s", f.list[aRow-1].Title)
			})
		} else {
			err = task.GetTaskCron().Remove(f.list[aRow-1].ID)
		}

		if err != nil {
			vcl.ShowMessage(fmt.Sprintf("任务状态更新失败: %s", err.Error()))
			return
		}
	}

	// 第六列 手动执行
	if aCol == 6 {
		func() {
			global.Logger.Sugar().Debugf("任务执行: %s", f.list[aRow-1].Title)
		}()
		return
	}

	// 第七列 查看日志
	if aCol == 7 {
		vcl.ShowMessage("查看日志")
		return
	}

	// 第八列 修改
	if aCol == 8 {
		ShowJobForm(f, f.list[aRow-1])
	}

	// 第九列 删除   最后一列，则表示是删除按钮
	if aCol == 9 {
		bdb.GetBdb().DeleteTask(f.list[aRow-1].ID)
	}

	f.GetData()
}

func (f *TFrm_job) GetData() {
	f.Table_data.Clear()
	f.Table_data.SetRowCount(2)

	var err error
	f.list, err = bdb.GetBdb().GetTasks(strings.TrimSpace(f.Edt_condition.Text()))
	if err != nil {
		vcl.ShowMessage(err.Error())
		return
	}

	length := int32(len(f.list))
	for i := int32(0); i < length; i++ {
		if f.list[i].Title == "" {
			continue
		}

		row := i + 1

		f.Table_data.SetCells(0, row, f.list[i].Title)
		f.Table_data.SetCells(1, row, f.list[i].Cron)
		f.Table_data.SetCells(2, row, strconv.Itoa(f.list[i].ScriptType))
		f.Table_data.SetCells(3, row, f.list[i].LastRunAt.String())
		f.Table_data.SetCells(4, row, f.list[i].Remarks)

		if f.list[i].Status == 1 {
			f.Table_data.SetCells(5, row, "停用")
		} else {
			f.Table_data.SetCells(5, row, "启用")
		}

		f.Table_data.SetCells(6, row, "执行")
		f.Table_data.SetCells(7, row, "查看日志")
		f.Table_data.SetCells(8, row, "修改")
		f.Table_data.SetCells(9, row, "删除")

		// 如果是最后一行，则不再添加新行
		if row == length {
			break
		}

		f.Table_data.SetRowCount(f.Table_data.RowCount() + 1)
	}
}
