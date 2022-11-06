<img align="right" width="150" src="https://user-images.githubusercontent.com/1475839/150521755-977f545b-4436-4059-87ac-1129541ad236.png" alt="donburi" title="donburi" /> <h1>Donburi</h1>

<img src="https://github.com/yohamta/donburi/actions/workflows/test.yaml/badge.svg" />

Donburi is just another Entity Component System library for Ebitengine inspired by [legion](https://github.com/amethyst/legion).

It aims to be a feature rich and high performance [ECS Library](https://en.wikipedia.org/wiki/Entity_component_system).

[GoDoc](https://pkg.go.dev/github.com/yohamta/donburi)

## Contents

- [Contents](#contents)
- [Features](#features)
- [Examples](#examples)
- [Installation](#installation)
- [Getting Started](#getting-started)
  - [Worlds](#worlds)
  - [Queries](#queries)
  - [Tags](#tags)
  - [Systems (Experimental)](#systems-experimental)
- [Features](#features-1)
  - [Math](#math)
  - [Transform](#transform)
  - [Events](#events)
- [Internal Design for `World`](#internal-design-for-world)
- [How to contribute?](#how-to-contribute)
- [Contributors](#contributors)

## Features

- It introduces the concept of [Archetype](https://docs.unity3d.com/Packages/com.unity.entities@0.2/manual/ecs_core.html), which allows us to query entities very efficiently based on the components layout.
- It is possible to combine `And`, `Or`, and `Not` conditions to perform complex queries for components.
- It avoids reflection on every frame for performance.
- Ability to dynamically add or remove components from an entity.
- APIs with Go Generics

There are many features that need to be added in the future (e.g., parent-child relationship, event-notification system etc).

## Examples

To check all examples, visit [this](https://github.com/yohamta/donburi/tree/master/examples) page.

The bunnymark example was adapted from [mizu](https://github.com/sedyh/mizu)'s code, which is made by [sedyh](https://github.com/sedyh). 

<a href="https://github.com/yohamta/donburi/tree/master/examples/bunnymark"> <img width="200" src="https://user-images.githubusercontent.com/1475839/150521292-9d3ec2c9-b96f-4cc1-a778-57dabfbd46b6.gif"></a> 

## Installation

```
go get github.com/yohamta/donburi
```

## Getting Started

### Worlds

```go
import "github.com/yohamta/donburi"

world := donburi.NewWorld()
```

Entities can be created via either `Create` (for a single entity) or `CreateMany` (for a collection of entities with the same component types). The world will create a unique ID for each entity upon insertion that you can use to refer to that entity later.

```go
// Component is any struct that holds some kind of data.
type PositionData struct {
  X, Y float64
}

type VelocityData struct {
  X, Y float64
}

// ComponentType represents kind of component which is used to create or query entities.
var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()

// Create an entity by specifying components that the entity will have.
// Component data will be initialized by default value of the struct.
entity = world.Create(Position, Velocity);

// You can use entity (it's a wrapper of int64) to get an Entry object from World
// which allows you to access the components that belong to the entity.
entry := world.Entry(entity)

position := donburi.Get[PositionData](entry, Position)
velocity := donburi.Get[VelocityData](entry, Velocity)
position.X += velocity.X
position.Y += velocity.y
```

We can define helper functions to get components for better readability. This was advice from [eliasdaler](https://github.com/eliasdaler).

```go
func GetPosition(entry *donburi.Entry) *PositionData {
  return donburi.Get[PositionData](entry, Position)
}

func GetVelocity(entry *donburi.Entry) *VelocityData {
  return donburi.Get[VelocityData](entry, Velocity)
}
```

Components can be added and removed through Entry objects.

```go
// Fetch the first entity with PlayerTag component
query := query.NewQuery(filter.Contains(PlayerTag))
// Query.FirstEntity() returns only the first entity that 
// matches the query.
if entry, ok := query.FirstEntity(world); ok {
  donburi.Add(entry, Position, &PositionData{
    X: 100,
    Y: 100,
  })
  donburi.Remove(entry, Velocity)
}
```

Entities can be removed from World with the World.Remove() as follows:

```go
if SomeLogic.IsDead(world, someEntity) {
  // World.Remove() removes the entity from the world.
  world.Remove(someEntity)
  // Deleted entities become invalid immediately.
  if world.Valid(someEntity) == false {
    println("this entity is invalid")
  }
}
```

### Queries

Queries allow for high performance and expressive iteration through the entities in a world, to get component references, test if an entity has a component or to add and remove components.

You can search for entities that have all of a set of components.

```go
// You can define a query by declaring what componet you want to find.
query := query.NewQuery(filter.Contains(Position, Velocity))

// You can then iterate through the entity found in the world
query.EachEntity(world, func(entry *donburi.Entry) {
  // An entry is an accessor to entity and its components.
  position := donburi.Get[PositionData](entry, Position)
  velocity := donburi.Get[VelocityData](entry, Velocity)
  
  position.X += velocity.X
  position.Y += velocity.Y
})
```

There are other types of filters such as `And`, `Or`, `Exact` and `Not`. You can combine them to find the target entities.

For example:

```go
// This query retrieves entities that have an NpcTag and no Position component.
query := query.NewQuery(filter.And(
  filter.Contains(NpcTag),
  filter.Not(filter.Contains(Position))))
```

If you need to determine if an entity has a component, there is `entry.HasComponent`

For example:

```go
// We have a query for all entities that have Position and Size, but also any of Sprite, Text or Shape.
query := query.NewQuery(
  filter.And(
    filter.Contains(Position, Size),
    filter.Or(
      filter.Contains(Sprite),
      filter.Contains(Text),
      filter.Contains(Shape),
    )
  )
)

// In our query we can check if the entity has some of the optional components before attempting to retrieve them
query.EachEntity(world, func(entry *donburi.Entry) {
  // We'll always be able to access Position and Size
  position := donburi.Get[PositionData](entry, Position)
  size := donburi.Get[SizeData](entry, Size)
  
  
  if entry.HasComponent(Sprite) {
    sprite := donburi.Get[SpriteData](entry, Sprite)
    // .. do sprite things
  }
  
  if entry.HasComponent(Text) {
    text := donburi.Get[TextData](entry, Text)
    // .. do text things
  }
  
  if entry.HasComponent(Shape) {
    shape := donburi.Get[ShapeData](entry, Shape)
    // .. do shape things
  }
  
})

```

### Tags

You can attach one or multiple "Tag" components to an entity. "Tag"s are just components with no data.

Here is the utility function to create a tag component.

```go
// This is the utility function to make tag component
func NewTag() *ComponentType {
  return NewComponentType(struct{}{})
}
```
Since "Tags" are components, they can be used in queries in the same way as components as follows:

```go
var EnemyTag = donburi.NewTag()
world.CreateMany(100, EnemyTag, Position, Velocity)

// Search entities with EnemyTag
query := query.NewQuery(filter.Contains(EnemyTag))
query.EachEntity(world, func(entry *donburi.Entry) {
  // Perform some operation on the Entities with the EnemyTag component.
}
```

### Systems (Experimental)

The [ECS package](https://github.com/yohamta/donburi/tree/main/ecs) provides so-called **System** feature in ECS which can be used together with a `World` instance.

See the [GoDoc](https://pkg.go.dev/github.com/yohamta/donburi/ecs) and [Example](https://github.com/yohamta/donburi/tree/master/examples/bunnymark_ecs).

How to create an ECS instance:

```go
import (
  "github.com/yohamta/donburi"
  ecslib "github.com/yohamta/donburi/ecs"
)

world := donburi.NewWorld()
ecs := ecslib.NewECS(world)
```

A `System` is created from just a function that receives an argument `(ecs *ecs.ECS)`.

```go
// Some System's function
func SomeFunction(ecs *ecs.ECS) {
  // ...
}

ecs.AddSystem(
  ecs.System{
    Update: SomeFunction,
  }
)
```

You can also provide `Draw()` functions for Systems. The `Layer` option allows you to control the order of calling `Draw()` functions and where to render. An `Layer` is just an integer value from the user's perspective. The default value of `Layer` is just `0`.

```go

const (
  LayerBackground ecslib.LayerID = iota
  LayerActors
  LayerFX
)

ecs.AddSystem(
  ecs.System{
    Layer:  LayerBackground,
    Update: UpdateBackground,
    Draw:   DrawBackground,
  }
)
```

Execute an ECS's `Update()` and `Draw()` to run systems as below:

```go
func (g *Game) Update() error {
  g.ecs.Update()
  return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
  screen.Clear()
  g.ecs.DrawLayer(LayerBackground, screen)
  g.ecs.DrawLayer(LayerBunnies, screen)
  g.ecs.DrawLayer(LayerMetrics, screen)
}
```

The `ecs.Create()` and `ecs.NewQuery()` wrapper-functions allow you to create and query entities on a certain layer:

```go
var layer0 ecs.LayerID = 0

// Create an entity on layer0
ecslib.Create(layer0, someComponents...)

// Create a query to iterate entities on layer0
queryForLayer0 := ecslib.NewQuery(layer0, filter.Contains(someComponent))
```

Example:

```go
import (
  ecslib "github.com/yohamta/donburi/ecs"
  "github.com/yohamta/donburi/features/hierarchy"
)

const (
  LayerBackground ecs.LayerID = iota
  LayerBunnies
  LayerMetrics
)

func newECS() *ecs.ECS {
  world := donburi.NewWorld()
  ecs := ecslib.NewECS(world)

  ecs.AddSystems(
    // Systems are executed in the order they are added.
    ecs.System{
      Update: system.NewSpawn().Update,
    },
    ecs.System{
      Layer: LayerBackground,
      Draw:  system.DrawBackground,
    },
    ecs.System{
      Layer:  LayerMetrics,
      Update: metrics.Update,
      Draw:   metrics.Draw,
    },
    ecs.System{
      Update: system.NewBounce(&g.bounds).Update,
    },
    ecs.System{
      Update: system.Velocity.Update,
    },
    ecs.System{
      Update: system.Gravity.Update,
    },
    ecs.System{
      Layer: LayerBunnies,
      Draw:  system.Render.Draw,
    },
  )

  return ecs
}
// ...

func (g *Game) Update() error {
  g.ecs.Update()
  return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
  screen.Clear()
  g.ecs.DrawLayer(LayerBackground, screen)
  g.ecs.DrawLayer(LayerBunnies, screen)
  g.ecs.DrawLayer(LayerMetrics, screen)
}

// ...

```


## Features

### Math

The [math package](https://github.com/yohamta/donburi/tree/main/features/math) provides the basic types (Vec2 etc) and helpers.

See the [GoDoc](https://pkg.go.dev/github.com/yohamta/donburi/features/math) for more details.

### Transform

The [transofrm package](https://github.com/yohamta/donburi/tree/main/features/transform) provides the `Tranform` Component and helpers.

`transform` package allows us to handle `position`, `rotation`, `scale` data relative to the parent.

This package was adapted from [ariplane](https://github.com/m110/airplanes)'s code, which is created by [m110](https://github.com/m110). 

Usage:
```go
w := donburi.NewWorld()

// setup parent
parent := w.Entry(w.Create(transform.Transform))

// set world position and scale for the parent
transform.SetWorldPosition(parent, dmath.Vec2{X: 1, Y: 2})
transform.SetWorldScale(parent, dmath.Vec2{X: 2, Y: 3})

// setup child
child := w.Entry(w.Create(transform.Transform))
donburi.SetValue(child, transform.Transform, transform.TransformData{
  LocalPosition: dmath.Vec2{X: 1, Y: 2},
  LocalRotation: 90,
  LocalScale:    dmath.Vec2{X: 2, Y: 3},
})

// add the child to the parent
transform.AppendChild(parent, child, false)

// get world position of the child with parent's position taken into account
pos := transform.WorldPosition(child)

// roatation
rot := transform.WorldRotation(child)

// scale
scale := transform.WorldScale(child)
```

How to remove chidren (= destroy entities):

```go
// Remove children
transform.RemoveChildrenRecursive(parent)

// Remove children and the parent
transform.RemoveRecursive(parent)
```

### Events

The [events package](https://pkg.go.dev/github.com/yohamta/donburi/features/events) allows us to send arbitrary data between systems in a Type-safe manner.

This package was adapted from [ariplane](https://github.com/m110/airplanes)'s code, which is created by [m110](https://github.com/m110). 

Example Usage:

```go

import "github.com/yohamta/donburi/features/events"

// Define any data
type EnemyKilled struct {
  EnemyID int
}

// Define an EventType with the type of the event data
var EnemyKilledEvent = events.NewEventType[EnemyKilled]()

// Create a world
world := donburi.NewWorld()

// Add handlers for the event
EnemyKilledEvent.Subscribe(world, LevelUp)
EnemyKilledEvent.Subscribe(world, UpdateScore)

// Sending an event
EnemyKilledEvent.Publish(world, EnemyKilled{EnemyID: 1})

// Process specific events
EnemyKilledEvent.ProcessEvents(world)

// Process all events
events.ProcessAllEvents(world)

// Receives the events
func LevelUp(w donburi.World, event EnemyKilled) {
  // do stuff
}

func UpdateScore(w donburi.World, event EnemyKilled) {
  // do stuff
}
```

## Internal Design for `World`

It is a bit complex, so please refer to the following diagram to understand it.

![arch](assets/architecture.png)

## How to contribute?

Feel free to contribute in any way you want. Share ideas, questions, submit issues, and create pull requests. Thanks!

## Contributors

<a href="https://github.com/yohamta/donburi/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=yohamta/donburi" />
</a>

Made with [contrib.rocks](https://contrib.rocks).
