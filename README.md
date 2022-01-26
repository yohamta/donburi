# <img align="right" width="150" src="https://user-images.githubusercontent.com/1475839/150521755-977f545b-4436-4059-87ac-1129541ad236.png" alt="donburi" title="donburi" /> Donburi

Donburi is just another Entity Component System library for Ebiten inspired by [legion](https://github.com/amethyst/legion).

It aims to be a feature rich and high performance [ECS Library](https://en.wikipedia.org/wiki/Entity_component_system).

## Contents

- [Features](#features)
- [Installation](#installation)
- [Examples](#examples)
- [Getting Started](#getting-started)

## Features

- It introduces the concept of [Archetype](https://docs.unity3d.com/Packages/com.unity.entities@0.2/manual/ecs_core.html), which allows us to query entities very efficiently based on the components layout.
- It is possible to combine And, Or, and Not conditions to perform complex queries for components.
- It avoids reflection on every frame and uses unsafe.Pointer for performance.
- Ability to dynamically add or remove components from an entity

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

Entities can be inserted via either `Create` (for a single entity) or `CreateMany` (for a collection of entities with the same component types). The world will create a unique ID for each entity upon insertion that you can use to refer to that entity later.

```go
// Component is any struct that holds some kind of data.
type PositionData struct {
  X, Y float64
}

type VelocityData struct {
  X, Y float64
}

// ComponentType represents kind of component which is used to create or query entities.
var Position = donburi.NewComponentType(PositionData{})
var Velocity = donburi.NewComponentType(VelocityData{})

// Create an entity by specifying components that the entity will have.
// Component data will be initialized by default value of the struct.
entity = world.Create(Position, Velocity);

// You can use entity (it's a wrapper of int64) to get an Entry object from World
// which allows you to access the components that belong to the entity.
entry := world.Entry(entity)

position := (*PositionData)(entry.Component(Position))
velocity := (*VelocityData)(entry.Component(Velocity))
position.X += velocity.X
position.Y += velocity.y
```

You can define helper functions to get components for better readability. This was advice from [eliasdaler](https://github.com/eliasdaler).

```go
func GetPosition(entry *donburi.Entry) *PositionData {
  return (*PositionData)(entry.Component(Position))
}

func GetVelocity(entry *donburi.Entry) *VelocityData {
  return (*VelocityData)(entry.Component(Velocity))
}
```

Components can be added and removed through Entry objects.

```go
// Fetch the first entity with PlayerTag component
query := query.NewQuery(filter.Contains(PlayerTag))
// Query.FirstEntity() returns only the first entity that 
// matches the query.
if entry, ok := query.FirstEntity(world); ok {
  entry.AddComponent(Position)
  entry.RemoveComponent(Velocity)
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

Queries allow for high performance and expressive iteration through the entities in a world, to find out what types of components are attached to it, to get component references, or to add and remove components.

You can search for entities which have all of a set of components.

```go
// You can define a query by declaring what componet you want to find.
query := query.NewQuery(filter.Contains(Position, Velocity))

// You can then iterate through the entity found in the world
query.EachEntity(world, func(entry *donburi.Entry) {
  // An entry is an accessor to entity and its components.
  var position *PositionData = (*PositionData)(entry.Component(Position))
  var velocity *VelocityData = (*VelocityData)(entry.Component(Velocity))
  
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

### Tags

You can attach "Tag" component to entity, which is just a component with no data.

Here is the utility function to create a tag component.

```go
// This is the utility function to make tag component
func NewTag() *ComponentType {
  return NewComponentType(struct{}{})
}
```
Since "Tags" are just components, they can be used in queries in the same way as components as follows:

```go
var EnemyTag = donburi.NewTag()
world.CreateMany(100, EnemyTag, Position, Velocity)

// Search entities with EnemyTag
query := query.NewQuery(filter.Contains(EnemyTag))
query.EachEntity(world, func(entry *donburi.Entry) {
  // Perform some operation on the Entities with the EnemyTag component.
}
```

### Systems

As of today, there is no function for the concept of "Systems" in ECS. It is assumed that operations are performed on entities using queries.
