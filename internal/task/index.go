package task

import (
	"fmt"
	"slices"
	"task/internal/bdb"

	cron "github.com/robfig/cron/v3"
)

var tc *TaskCron

func GetTaskCron() *TaskCron {
	if tc == nil {
		tc = NewTaskCron()
	}
	return tc
}

func InitTask() {
	GetTaskCron()
	tc.Start()
}

type TaskCron struct {
	*cron.Cron
	list []*TaskInfo
}

func (tc *TaskCron) GetTask(id int) *TaskInfo {
	for _, info := range tc.list {
		if info.Task.ID == id {
			return info
		}
	}

	return nil
}

type TaskInfo struct {
	ID cron.EntryID `json:"cronId"` // 定时任务ID
	*bdb.Task
}

type TaskRun struct {
	ID   int
	Path string
	Job  func()
}



func (tr TaskRun) Run() {
	// 每次执行任务前，先记录日志
	fmt.Printf("TaskRun: %d\n", tr.ID)
	tr.Job()
}

func NewTaskRun(id int, job func()) TaskRun {
	return TaskRun{
		ID:  id,
		Job: job,
	}
}

func NewTaskCron() *TaskCron {

	c := cron.New(cron.WithSeconds())

	return &TaskCron{
		Cron: c,
		list: make([]*TaskInfo, 0),
	}
}

func (tc *TaskCron) Start() {
	tc.Cron.Start()
}

func (tc *TaskCron) AddFunc(info *bdb.Task, cmd func()) error {
	tRun := NewTaskRun(info.ID, cmd)
	entryID, err := tc.Cron.AddJob(info.Cron, tRun)
	if err != nil {
		return err
	}

	tc.list = append(tc.list, &TaskInfo{
		ID:   entryID,
		Task: info,
	})

	return nil
}

func (tc *TaskCron) Remove(id int) error {
	job := tc.GetTask(id)
	if job == nil {
		return fmt.Errorf("[%d] task not found", id)
	}

	tc.list = slices.DeleteFunc(tc.list, func(info *TaskInfo) bool {
		return info.Task.ID == id
	})

	tc.Cron.Remove(job.ID)
	return nil
}

func (tc *TaskCron) CronList() []*TaskInfo {
	return tc.list
}
