package repository

import (
	"clase02-mongo/internal/dao"
	"clase02-mongo/internal/domain"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

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
func (r *MongoItemsRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	itemDAO := dao.FromDomain(item) // Convertir a modelo DAO

	itemDAO.ID = primitive.NewObjectID()
	itemDAO.CreatedAt = time.Now().UTC()
	itemDAO.UpdatedAt = time.Now().UTC()

	// Insertar en DB
	_, err := r.col.InsertOne(ctx, itemDAO)
	if err != nil {
		// Podemos manejar errores específicos de MongoDB, como claves duplicadas
		// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#IsDuplicateKeyError
		// Esto es útil si tenemos restricciones de unicidad en la colección
		if mongo.IsDuplicateKeyError(err) {
			return domain.Item{}, errors.New("item with the same ID already exists")
		}

		// Error genérico
		return domain.Item{}, err
	}

	return itemDAO.ToDomain(), nil // Convertir de vuelta a Domain
}

// GetByID busca un item por su ID
// Consigna 2: Validar que el ID sea un ObjectID válido
func (r *MongoItemsRepository) GetByID(ctx context.Context, id string) (domain.Item, error) {
	// Validar que el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Item{}, errors.New("invalid ID format")
	}

	// Buscar en DB
	var itemDAO dao.Item
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&itemDAO)
	if err != nil {
		// Manejar caso de no encontrado
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Item{}, errors.New("item not found")
		}
		return domain.Item{}, err
	}

	return itemDAO.ToDomain(), nil
}

// Update actualiza un item existente
// Consigna 3: Update parcial + actualizar updatedAt
func (r *MongoItemsRepository) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {
	// Validar que el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Item{}, errors.New("invalid ID format")
	}

	// Preparar los campos a actualizar
	update := bson.M{
		"$set": bson.M{
			"name":       item.Name,
			"price":      item.Price,
			"updated_at": time.Now().UTC(), // Actualizar timestamp
		},
	}

	// Ejecutar la actualización
	result, err := r.col.UpdateByID(ctx, objID, update)
	if err != nil {
		return domain.Item{}, err
	}
	if result.MatchedCount == 0 {
		return domain.Item{}, errors.New("item not found")
	}

	// Retornar el item actualizado
	return r.GetByID(ctx, id)
}

// Delete elimina un item por ID
// Consigna 4: Eliminar documento de DB
func (r *MongoItemsRepository) Delete(ctx context.Context, id string) error {
	// Validar que el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	// Ejecutar la eliminación
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("item not found")
	}

	return nil
}
