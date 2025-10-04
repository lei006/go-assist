package taskgroup

type TaskGroup struct {
	TaskCount int64 //任务数量

	taskChan chan TaskGroupData
}

func MakeNewTaskGroup() *TaskGroup {
	task_group := TaskGroup{}
	task_group.taskChan = make(chan TaskGroupData)
	return &task_group
}

func (task_grop *TaskGroup) SetTaskCount(count int64) {
	if task_grop.TaskCount < count {
		count := count - task_grop.TaskCount
		for i := int64(0); i < count; i++ {
			new_task := MakeNewTaskUnit(i, &task_grop.taskChan)
			new_task.Run()
		}
	} else {
		//发送指定数量的退出信号
		count := task_grop.TaskCount - count
		for i := int64(0); i < count; i++ {
			task_grop.ExitOneTask()
		}
	}
}

func (task_grop *TaskGroup) ExitOneTask() {

	newTaskGroupData := TaskGroupData{
		IsExit: true,
	}
	task_grop.taskChan <- newTaskGroupData

}

func (task_grop *TaskGroup) Push(data TaskData) {

	newTaskGroupData := TaskGroupData{
		Data: data,
	}
	task_grop.taskChan <- newTaskGroupData

}
