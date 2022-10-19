package ecs

import "testing"

func TestLayer(t *testing.T) {
	l := getLayer(0)
	ll := getLayer(1)

	if l.id != 0 {
		t.Errorf("layer id is not 0")
	}
	if ll.id != 1 {
		t.Errorf("layer id is not 1")
	}
	if l.tag == ll.tag {
		t.Errorf("layer tag is same")
	}
	if l.tag.Name() != "Layer0" {
		t.Errorf("layer tag name is not Layer0")
	}
}
