![](./assets/images/example.gif)

# Resolv Example using ECS architecture

This example code is created based on the example of [resolv](https://github.com/SolarLune/resolv) created by [SolarLune](https://github.com/SolarLune).

## How to play

Arrow key to move, X to jump.

## Architecture

This example code is based on the ECS architecture using [donburi](https://github.com/yohamta/donburi).

The ECS structure is as follows:

* `scenes` - contains the main logic for setup and execute ECS systems.
* `components` - contains all components. **Components** are Go structs with no behavior other than occasional helper methods.
* `systems` - contains all systems. **Systems** keep the logic of the game. Each system works on entities with a specific set of components.
* `archetypes` - contains helper functions for creating entities with specific sets of components.
* `factory` - contains helper functions for creating and initializing data of components. `factory` uses `archetype` helpers.
