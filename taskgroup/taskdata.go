package taskgroup

type TaskData interface {
	Run()
	GetID() int64
}

type TaskGroupData struct {
	IsExit bool
	Data   TaskData
}
