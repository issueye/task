package home

import (
	"task/internal/bdb"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

func ShowJobForm(owner vcl.IComponent, data *bdb.Task) {
	if Frm_job_form == nil {
		Frm_job_form = NewFrm_job_form(owner)
	}

	if data != nil {
		Frm_job_form.SetData(data)
	} else {
		Frm_job_form.Reset()
	}

	Frm_job.SetPosition(types.PoMainFormCenter)
	Frm_job.SetFormStyle(types.FsStayOnTop)

	Frm_job_form.ShowModal()
}

// ::private::
type TFrm_job_formFields struct {
	data *bdb.Task
}

func (f *TFrm_job_form) OnFormCreate(sender vcl.IObject) {

}

func (f *TFrm_job_form) SetEvents() {
	f.Btn_save.SetOnClick(f.OnBtn_saveClick)
	f.Btn_cancel.SetOnClick(f.OnBtn_cancelClick)
}

func (f *TFrm_job_form) Reset() {
	f.Edt_title.SetText("")
	f.Edt_cron.SetText("")
	f.Edt_script_path.SetText("")
	f.Mmo_remarks.SetText("")
	f.RB_01.SetChecked(true)
}

func (f *TFrm_job_form) SetData(data *bdb.Task) {
	f.data = data
	f.Edt_title.SetText(data.Title)
	f.Edt_cron.SetText(data.Cron)
	f.Edt_script_path.SetText(data.ScritpPath)
	f.Mmo_remarks.SetText(data.Remarks)

	if data.ScriptType == 0 {
		f.RB_01.SetChecked(true)
	}

	if data.ScriptType == 1 {
		f.RB_02.SetChecked(true)
	}
	if data.ScriptType == 2 {
		f.RB_03.SetChecked(true)
	}
}

func (f *TFrm_job_form) OnBtn_saveClick(sender vcl.IObject) {
	if f.data == nil {
		f.SaveData()
	} else {
		f.UpdateData()
	}

	f.Close()
}

func (f *TFrm_job_form) OnBtn_cancelClick(sender vcl.IObject) {
	f.Close()
}

func (f *TFrm_job_form) UpdateData() {
	if f.data == nil {
		return
	}

	f.data.Title = f.Edt_title.Text()
	f.data.Cron = f.Edt_cron.Text()
	f.data.ScriptType = f.GetSelect()
	f.data.ScritpPath = f.Edt_script_path.Text()
	f.data.Remarks = f.Mmo_remarks.Text()

	err := bdb.GetBdb().UpdateTask(f.data)
	if err != nil {
		vcl.ShowMessage(err.Error())
		return
	}
}

func (f *TFrm_job_form) SaveData() {
	data := new(bdb.Task)
	data.Title = f.Edt_title.Text()
	data.Cron = f.Edt_cron.Text()
	data.ScriptType = f.GetSelect()
	data.ScritpPath = f.Edt_script_path.Text()
	data.Remarks = f.Mmo_remarks.Text()

	err := bdb.GetBdb().CreateTask(data)
	if err != nil {
		vcl.ShowMessage(err.Error())
		return
	}
}

func (f *TFrm_job_form) GetSelect() int {
	if f.RB_02.Checked() {
		return 1
	}

	if f.RB_03.Checked() {
		return 2
	}

	return 0
}
