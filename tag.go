package donburi

type Tag string

// NewTag is an utility to create a tag component.
// Which is just an component that contains no data.
func NewTag(name string) *ComponentType[Tag] {
	c := NewComponentType[Tag](Tag(name))
	c.SetName(name)
	return c
}
