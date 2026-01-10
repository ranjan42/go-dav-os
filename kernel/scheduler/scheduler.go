package scheduler

import "unsafe"

const StackSize = 4096
const MaxTasks = 16

type TaskState int

const (
	TaskRunnable TaskState = iota
	TaskRunning
	TaskWaiting
)

type Task struct {
	ID    int
	ESP   uint32
	State TaskState
	Stack [StackSize]byte
}

var (
	tasks       [MaxTasks]*Task
	taskCount   int
	currentTask *Task
	nextID      int = 1

	// Static allocation for tasks to avoid 'newobject' heap allocation
	taskPool [MaxTasks]Task
)

// CpuSwitch is defined in switch.s
func CpuSwitch(oldESP *uint32, newESP uint32)

func Init() {
	taskCount = 0
	// Init initial task (0)
	t := &taskPool[0]
	t.ID = 0
	t.State = TaskRunning

	tasks[0] = t
	taskCount = 1
	currentTask = t
}

func NewTask(entry func()) *Task {
	if taskCount >= MaxTasks {
		return nil
	}

	idx := taskCount
	t := &taskPool[idx]
	t.ID = nextID
	nextID++
	t.State = TaskRunnable

	// Setup stack
	// Stack grows down from t.Stack[StackSize]

	sp := uintptr(unsafe.Pointer(&t.Stack[StackSize-1]))
	sp = sp & ^uintptr(15)

	// Store return addr (entry) and register placeholders
	sp -= 4
	*(*uintptr)(unsafe.Pointer(sp)) = *(*uintptr)(unsafe.Pointer(&entry))

	sp -= 16 // 4 * 4 bytes for EBP, EBX, ESI, EDI

	t.ESP = uint32(sp)

	tasks[idx] = t
	taskCount++
	return t
}

func Schedule() {
	if taskCount <= 1 {
		return
	}

	oldTask := currentTask

	nextIndex := -1
	currentIndex := -1

	for i := 0; i < taskCount; i++ {
		if tasks[i] == currentTask {
			currentIndex = i
			break
		}
	}

	// Round robin
	for i := 1; i < taskCount; i++ {
		idx := (currentIndex + i) % taskCount
		if tasks[idx].State == TaskRunnable {
			nextIndex = idx
			break
		}
	}

	if nextIndex == -1 || tasks[nextIndex] == currentTask {
		return
	}

	newTask := tasks[nextIndex]

	oldTask.State = TaskRunnable
	newTask.State = TaskRunning
	currentTask = newTask

	CpuSwitch(&oldTask.ESP, newTask.ESP)
}

func CurrentTaskID() int {
	if currentTask == nil {
		return -1
	}
	return currentTask.ID
}
