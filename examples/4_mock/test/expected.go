package test

type HelloMock struct {
	HelloFunc      func() string
	SayHelloToFunc func(param0 string) string
}

func (m *HelloMock) Hello() string {
	return m.HelloFunc()
}
func (m *HelloMock) SayHelloTo(param0 string) string {
	return m.SayHelloToFunc(param0)
}
