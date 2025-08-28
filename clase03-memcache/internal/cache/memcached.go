package cache

// TODO(Clase):
// - Implementar un wrapper de Memcached similar al de la solución
// - API sugerida:
//   type Client struct { ... }
//   func New(addr string, ttl time.Duration) *Client
//   func (c *Client) GetJSON(key string, dest any) (bool, error)
//   func (c *Client) SetJSON(key string, v any) error
//   func (c *Client) Delete(key string)
//   func (c *Client) SelfTest(ctx context.Context) error
