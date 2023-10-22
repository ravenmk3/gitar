package fslock

type fsLock struct {
	filename string
}

func New(filename string) Lock {
	return &fsLock{filename: filename}
}

func (l *fsLock) Lock() error {
	return nil
}

func (l *fsLock) TryLock() error {
	return nil
}

func (l *fsLock) Unlock() error {
	return nil
}
