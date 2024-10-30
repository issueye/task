package home

import (
	"context"
	"task/internal/global"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ying32/govcl/vcl"
)

// ::private::
type TFrm_taskFields struct {
	ctx     context.Context
	cancel  context.CancelFunc
	FrmMain *TFrm_main
}

func (f *TFrm_task) OnFormCreate(sender vcl.IObject) {
	f.ctx, f.cancel = context.WithCancel(context.Background())
	f.Monitor(f.ctx)
	f.SetEvents()
}

func (f *TFrm_task) SetEvents() {
	f.SetOnDestroy(f.OnFormDestroy)
	f.SetOnShow(f.OnFormShow)
}

func (f *TFrm_task) OnFormShow(sender vcl.IObject) {
	f.SetVisible(false)
}

func (f *TFrm_task) OnFormDestroy(sender vcl.IObject) {
	f.cancel()
}

func (f *TFrm_task) Monitor(ctx context.Context) {
	showMsg, err := global.PubSub.Subscribe(context.Background(), global.TOPIC_SHOW_HOME)
	if err != nil {
		return
	}

	go func(c context.Context) {
		for {
			select {
			case msg := <-showMsg:
				f.ShowForm(msg)
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}

func (f *TFrm_task) ShowForm(msg *message.Message) {
	msg.Ack()
	global.Logger.Sugar().Debugf("show home")
	vcl.ThreadSync(func() {
		ShowHomeForm(nil)
	})
}
