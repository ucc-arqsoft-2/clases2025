package cache

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"example.com/gin-memcached-base/internal/models"
)

// MemcachedClient implementa la interfaz Cache usando memcached real
type MemcachedClient struct {
	addr    string
	timeout time.Duration
}

// NewMemcachedClient crea una nueva conexión a memcached
func NewMemcachedClient(addr string) *MemcachedClient {
	return &MemcachedClient{
		addr:    addr,
		timeout: 5 * time.Second,
	}
}

// Set guarda un item en memcached
func (m *MemcachedClient) Set(key string, item models.Item) error {
	conn, err := net.DialTimeout("tcp", m.addr, m.timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to memcached: %w", err)
	}
	defer conn.Close()

	// Serializar el item a JSON
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	// Comando SET de memcached: set <key> <flags> <exptime> <bytes>\r\n<data>\r\n
	command := fmt.Sprintf("set %s 0 3600 %d\r\n%s\r\n", key, len(data), data)
	
	_, err = conn.Write([]byte(command))
	if err != nil {
		return fmt.Errorf("failed to write to memcached: %w", err)
	}

	// Leer respuesta
	reader := bufio.NewReader(conn)
	response, _, err := reader.ReadLine()
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if !strings.HasPrefix(string(response), "STORED") {
		return fmt.Errorf("memcached set failed: %s", string(response))
	}

	return nil
}

// Get obtiene un item de memcached
func (m *MemcachedClient) Get(key string) (models.Item, error) {
	conn, err := net.DialTimeout("tcp", m.addr, m.timeout)
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to connect to memcached: %w", err)
	}
	defer conn.Close()

	// Comando GET de memcached: get <key>\r\n
	command := fmt.Sprintf("get %s\r\n", key)
	
	_, err = conn.Write([]byte(command))
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to write to memcached: %w", err)
	}

	// Leer respuesta
	reader := bufio.NewReader(conn)
	
	// Primera línea: VALUE <key> <flags> <bytes>\r\n o END\r\n
	line, err := reader.ReadString('\n')
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to read response: %w", err)
	}

	line = strings.TrimSpace(line)
	if line == "END" {
		return models.Item{}, fmt.Errorf("item no encontrado")
	}

	// Parsear VALUE line
	parts := strings.Split(line, " ")
	if len(parts) < 4 || parts[0] != "VALUE" {
		return models.Item{}, fmt.Errorf("invalid response format: %s", line)
	}

	// Leer los datos JSON
	dataLine, err := reader.ReadString('\n')
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to read data: %w", err)
	}

	// Deserializar JSON
	var item models.Item
	err = json.Unmarshal([]byte(strings.TrimSpace(dataLine)), &item)
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	// Leer línea END
	reader.ReadString('\n')

	return item, nil
}