package ecs

// UpdateSystem is a system that updates the world.
type System func(ecs *ECS)

// DrawSystem is a system that draws the world.
type RendererWithArg[T Arg] func(ecs *ECS, arg *T)

// Arg is an argument of the renderer.
type Arg interface{}
