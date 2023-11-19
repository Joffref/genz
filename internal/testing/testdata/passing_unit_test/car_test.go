package main

import "testing"

func Test_WorkingTest(t *testing.T) {
	t.Log("This test should pass")
	c := &Car{model: "Ford"}
	if c.GetModel() != "Ford" {
		t.Errorf("Expected model to be Ford, got %s", c.GetModel())
	}
}
