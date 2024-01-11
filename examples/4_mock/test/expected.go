package test

type HelloMock struct {
	SayHelloToFunc func(name string) string
	HelloFunc      func() string
}

func (m *HelloMock) SayHelloTo(name string) string {
	return m.SayHelloToFunc(name)
}

func (m *HelloMock) Hello() string {
	return m.HelloFunc()
}
