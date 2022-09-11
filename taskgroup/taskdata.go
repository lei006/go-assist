package taskgroup

type TaskData interface {
	Run(int64)
}

type TaskGroupData struct {
	IsExit bool
	Data   TaskData
}
