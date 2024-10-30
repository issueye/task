package home

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/win"
)

func ShowHomeForm(owner vcl.IComponent) {
	if Frm_main == nil {
		Frm_main = NewFrm_main(owner)
	}

	Frm_main.Show()

	Frm_main.SetPosition(types.PoMainFormCenter)
	Frm_main.SetFormStyle(types.FsStayOnTop)

	Frm_main.Update()

	width := vcl.Screen.Width()
	height := vcl.Screen.Height()
	Frm_main.SetLeft((width - Frm_main.Width()) / 2)
	Frm_main.SetTop((height - Frm_main.Height()) / 2)

	win.SetWindowPos(Frm_main.Handle(), win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)

}

// ::private::
type TFrm_mainFields struct {
}

func (f *TFrm_main) OnFormCreate(sender vcl.IObject) {
	f.TForm.SetOnShow(f.OnFormShow)
}

func (f *TFrm_main) OnFormShow(sender vcl.IObject) {
}
