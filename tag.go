package donburi

// NewTag is an utility to create a tag component.
// Which is just an component that contains no data.
func NewTag() *ComponentType[struct{}] {
	return NewComponentType[struct{}]()
}
