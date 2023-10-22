package fslock

type Lock interface {
	Lock() error
	TryLock() error
	Unlock() error
}
