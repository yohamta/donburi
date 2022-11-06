package donburi

import "github.com/yohamta/donburi/internal/component"

// NewTag is an utility to create a tag component.
// Which is just an component that contains no data.
func NewTag() *component.ComponentType[struct{}] {
	return NewComponentType[struct{}]()
}
