package test

func (v Car) GetModel() string {
	return v.Model
}

func (v Car) GetWheels() []Wheel {
	return v.Wheels
}
