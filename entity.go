package donburi

import "github.com/yohamta/donburi/internal/storage"

// Entity is identifier of an entity.
// Entity is just a wrapper of uint64.
type Entity = storage.Entity

// Null represents a invalid entity which is zero.
var Null = storage.Null
