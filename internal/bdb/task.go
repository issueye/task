package bdb

import (
	"encoding/json"
	"strings"
	"task/internal/global"
	"time"

	"go.etcd.io/bbolt"
)

var TaskName = []byte("task")

type Task struct {
	ID         int       `json:"id"`          // 任务ID
	Title      string    `json:"title"`       // 任务名称
	Cron       string    `json:"cron"`        // cron表达式
	Status     int       `json:"status"`      // 任务状态
	ScriptType int       `json:"script_type"` // 脚本类型
	ScritpPath string    `json:"script_path"` // 脚本路径
	Remarks    string    `json:"remarks"`     // 备注
	LastRunAt  time.Time `json:"last_run_at"` // 上次运行时间
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`  // 更新时间
}

func (t *Task) TableName() string {
	return "task"
}

func (b *Bdb) CreateTask(data *Task) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(TaskName)
		if err != nil {
			return err
		}

		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}

		global.Logger.Sugar().Debugf("新增数据， id : %d", id)

		data.ID = int(id)

		buf, err := json.Marshal(data)
		if err != nil {
			return err
		}

		return bucket.Put(itob(data.ID), buf)
	})
}

func (b *Bdb) GetById(id int) (*Task, error) {
	var task Task
	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(TaskName)
		if bucket == nil {
			return ErrBucketNotFound
		}

		data := bucket.Get(itob(id))
		if data == nil {
			return ErrRecordNotFound
		}

		return json.Unmarshal(data, &task)
	})

	return &task, err
}

func (b *Bdb) UpdateTask(data *Task) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(TaskName)
		if bucket == nil {
			return ErrBucketNotFound
		}

		buf, err := json.Marshal(data)
		if err != nil {
			return err
		}

		return bucket.Put(itob(data.ID), buf)
	})
}

func (b *Bdb) DeleteTask(id int) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(TaskName)
		if bucket == nil {
			return ErrBucketNotFound
		}

		return bucket.Delete(itob(id))
	})
}

func (b *Bdb) GetTasks(condition string) ([]*Task, error) {
	var tasks []*Task
	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(TaskName)
		if bucket == nil {
			return ErrBucketNotFound
		}

		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var task Task
			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}

			if condition != "" {
				if strings.Contains(task.Title, condition) || strings.Contains(task.Remarks, condition) {
					tasks = append(tasks, &task)
					continue
				}
			}

			tasks = append(tasks, &task)
		}
		return nil
	})

	return tasks, err
}
