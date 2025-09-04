package cache

import "time"

// Cache define la interfaz para nuestro wrapper
type Cache interface {
    Get(key string) ([]byte, bool)
    Set(key string, val []byte, ttl time.Duration) error
    Delete(key string) error
    // TODO: devolver la lista de claves cacheadas para poder inspeccionar desde la API
    Keys() ([]string, error)
}
