package mymath_test

import (
	"testing"
	"mymath"
	"log"
)

func TestAdd(t *testing.T) {
	ret := mymath.Add(2, 3)
	if ret != 5 {
		t.Error("Expected 5, got ", ret)
	}
	log.Print("this my test!")
}
