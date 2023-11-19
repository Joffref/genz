package test

func (c *Car) GetModel() string {
	return c.model
}

func (c *Car) GetWheels() []Wheel {
	return c.wheels
}
