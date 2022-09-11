package taskgroup

var Default *TaskGroup

func init() {
	Default = MakeNewTaskGroup()
}
