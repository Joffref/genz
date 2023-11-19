package main

import "testing"

func Test_FailingTest(t *testing.T) {
	t.Errorf("This test should fail")
}
