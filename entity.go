package donburi

import "github.com/yohamta/donburi/internal/entity"

// Entity is identifier of an entity.
// Entity is just a wrapper of uint64.
type Entity = entity.Entity

// Null represents a invalid entity which is zero.
var Null = entity.Null
