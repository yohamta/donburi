package donburi

type Tag string

// NewTag is an utility to create a tag component.
// Which is just an component that contains no data.
// Specify a string as the first and only parameter if you wish to name the component.
func NewTag(opts ...any) *ComponentType[Tag] {
	if len(opts) == 0 {
		return NewComponentType[Tag]()
	}
	first, ok := opts[0].(string)
	if !ok {
		return NewComponentType[Tag]()
	}
	c := NewComponentType[Tag](Tag(first))
	c.SetName(first)
	return c
}
