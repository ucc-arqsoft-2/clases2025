package repository

import (
	"clase02-mongo/internal/dao"
	"clase02-mongo/internal/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ItemsRepository define las operaciones de datos para Items
// Patrón Repository: abstrae el acceso a datos del resto de la aplicación
type ItemsRepository interface {
	// List retorna todos los items de la base de datos
	List(ctx context.Context) ([]domain.Item, error)

	// Create inserta un nuevo item en MongoDB
	Create(ctx context.Context, item domain.Item) (domain.Item, error)

	// GetByID busca un item por su ID
	GetByID(ctx context.Context, id string) (domain.Item, error)

	// Update actualiza un item existente
	Update(ctx context.Context, id string, item domain.Item) (domain.Item, error)

	// Delete elimina un item por ID
	Delete(ctx context.Context, id string) error
}

// MongoItemsRepository implementa ItemsRepository usando MongoDB
type MongoItemsRepository struct {
	col *mongo.Collection // Referencia a la colección "items" en MongoDB
}

// NewMongoItemsRepository crea una nueva instancia del repository
// Recibe una referencia a la base de datos MongoDB
func NewMongoItemsRepository(db *mongo.Database) ItemsRepository {
	return &MongoItemsRepository{
		col: db.Collection("items"), // Conecta con la colección "items"
	}
}

// List obtiene todos los items de MongoDB
func (r *MongoItemsRepository) List(ctx context.Context) ([]domain.Item, error) {
	// ⏰ Timeout para evitar que la operación se cuelgue
	// Esto es importante en producción para no bloquear indefinidamente
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 🔍 Find() sin filtros retorna todos los documentos de la colección
	// bson.M{} es un filtro vacío (equivale a {} en MongoDB shell)
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx) // ⚠️ IMPORTANTE: Siempre cerrar el cursor para liberar recursos

	// 📦 Decodificar resultados en slice de DAO (modelo MongoDB)
	// Usamos el modelo DAO porque maneja ObjectID y tags BSON
	var daoItems []dao.Item
	if err := cur.All(ctx, &daoItems); err != nil {
		return nil, err
	}

	// 🔄 Convertir de DAO a Domain (para la capa de negocio)
	// Separamos los modelos: DAO para MongoDB, Domain para lógica de negocio
	domainItems := make([]domain.Item, len(daoItems))
	for i, daoItem := range daoItems {
		domainItems[i] = daoItem.ToDomain() // Función definida en dao/Item.go
	}

	return domainItems, nil
}

// Create inserta un nuevo item en MongoDB
// Consigna 1: Validar name y price >= 0, agregar timestamps
func (r *MongoItemsRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	return domain.Item{}, errors.New("TODO: implementar Create")
}

// GetByID busca un item por su ID
// Consigna 2: Validar que el ID sea un ObjectID válido
func (r *MongoItemsRepository) GetByID(ctx context.Context, id string) (domain.Item, error) {
	return domain.Item{}, errors.New("TODO: implementar GetByID")
}

// Update actualiza un item existente
// Consigna 3: Update parcial + actualizar updatedAt
func (r *MongoItemsRepository) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {
	return domain.Item{}, errors.New("TODO: implementar Update")
}

// Delete elimina un item por ID
// Consigna 4: Eliminar documento de MongoDB
func (r *MongoItemsRepository) Delete(ctx context.Context, id string) error {
	return errors.New("TODO: implementar Delete")
}
