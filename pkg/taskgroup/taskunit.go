package taskgroup

// 任务调度器
type TaskUnit struct {
	taskChan *chan TaskGroupData
	Id       int64
}

func MakeNewTaskUnit(id int64, ch *chan TaskGroupData) *TaskUnit {
	tmp := &TaskUnit{
		taskChan: ch,
		Id:       id,
	}

	return tmp
}

func (task *TaskUnit) Run() {
	go func() {

		for {
			task_data := <-*task.taskChan
			if task_data.IsExit {
				//立刻退出..
				break
			}
			task_data.Data.Run(task.Id)

		}

	}()

}
