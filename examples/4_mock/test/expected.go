package test

type HelloMock struct {
	SayHelloToFunc func(param0 string) string
}

func (m *HelloMock) SayHelloTo(param0 string) string {
	return m.SayHelloToFunc(param0)
}
