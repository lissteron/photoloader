package elog

func NewStub() Logger {
	logger, _ := New("debug")

	return logger
}
