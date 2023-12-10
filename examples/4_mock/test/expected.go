package test

type HelloMock struct {
	SayHelloToFunc func(param0 string) string
	HelloFunc      func() string
}

func (m *HelloMock) SayHelloTo(param0 string) string {
	return m.SayHelloToFunc(param0)
}

func (m *HelloMock) Hello() string {
	return m.HelloFunc()
}
