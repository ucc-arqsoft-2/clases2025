package repository

import (
	"clase02-mongo/internal/dao"
	"clase02-mongo/internal/domain"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoItemsRepository implementa ItemsRepository usando DB
type MongoItemsRepository struct {
	col *mongo.Collection // Referencia a la colección "items" en DB
}

// NewMongoItemsRepository crea una nueva instancia del repository
// Recibe una referencia a la base de datos DB
func NewMongoItemsRepository(ctx context.Context, uri, dbName, collectionName string) *MongoItemsRepository {
	opt := options.Client().ApplyURI(uri)
	opt.SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
		return nil
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		log.Fatalf("Error pinging DB: %v", err)
		return nil
	}

	return &MongoItemsRepository{
		col: client.Database(dbName).Collection(collectionName), // Conecta con la colección "items"
	}
}

// List obtiene todos los items de DB
func (r *MongoItemsRepository) List(ctx context.Context) ([]domain.Item, error) {
	// ⏰ Timeout para evitar que la operación se cuelgue
	// Esto es importante en producción para no bloquear indefinidamente
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 🔍 Find() sin filtros retorna todos los documentos de la colección
	// bson.M{} es un filtro vacío (equivale a {} en DB shell)
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx) // ⚠️ IMPORTANTE: Siempre cerrar el cursor para liberar recursos

	// 📦 Decodificar resultados en slice de DAO (modelo DB)
	// Usamos el modelo DAO porque maneja ObjectID y tags BSON
	var daoItems []dao.Item
	if err := cur.All(ctx, &daoItems); err != nil {
		return nil, err
	}

	// 🔄 Convertir de DAO a Domain (para la capa de negocio)
	// Separamos los modelos: DAO para DB, Domain para lógica de negocio
	domainItems := make([]domain.Item, len(daoItems))
	for i, daoItem := range daoItems {
		domainItems[i] = daoItem.ToDomain() // Función definida en dao/Item.go
	}

	return domainItems, nil
}

// Create inserta un nuevo item en DB
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
// Consigna 4: Eliminar documento de DB
func (r *MongoItemsRepository) Delete(ctx context.Context, id string) error {
	return errors.New("TODO: implementar Delete")
}
