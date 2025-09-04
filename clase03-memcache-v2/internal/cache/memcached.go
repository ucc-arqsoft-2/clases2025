package cache

import (
    "time"

    "github.com/bradfitz/gomemcache/memcache"
)

// NOTA: esta implementación es mínima a propósito para la clase.
// Ejercicios:
// 1) Agregar un 'índice' de claves (por ejemplo guardado bajo 'cache:index' en JSON)
// 2) Implementar Keys() leyendo ese índice
// 3) En Delete(key), remover la clave del índice si existe

type Memcached struct {
    c *memcache.Client
}

func NewMemcached(addr string) *Memcached {
    return &Memcached{c: memcache.New(addr)}
}

func (m *Memcached) Get(key string) ([]byte, bool) {
    it, err := m.c.Get(key)
    if err != nil {
        return nil, false
    }
    return it.Value, true
}

func (m *Memcached) Set(key string, val []byte, ttl time.Duration) error {
    seconds := int32(ttl.Seconds())
    if seconds <= 0 { seconds = 60 }
    return m.c.Set(&memcache.Item{Key: key, Value: val, Expiration: seconds})
}

func (m *Memcached) Delete(key string) error {
    return m.c.Delete(key)
}

func (m *Memcached) Keys() ([]string, error) {
    // TODO: devolver claves reales usando el índice
    return []string{}, nil
}
