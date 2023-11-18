package main

import "testing"

func TestCarGetModel(t *testing.T) {
	c := &Car{}
	if c.GetModel() != c.model {
		t.Errorf("Expected %s, got %s", c.model, c.GetModel())
	}
	c.model = "Ford"
	if c.GetModel() != "Ford" {
		t.Errorf("Expected %s, got %s", c.model, c.GetModel())
	}
	c.model = "Ferrari"
	if c.GetModel() != "Ferrari" {
		t.Errorf("Expected %s, got %s", c.model, c.GetModel())
	}
}
