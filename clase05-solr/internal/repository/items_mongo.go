package repository

import (
	"clase05-solr/internal/dao"
	"clase05-solr/internal/domain"
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
func (r *MongoItemsRepository) List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error) {
	// ⏰ Timeout para evitar que la operación se cuelgue
	// Esto es importante en producción para no bloquear indefinidamente
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 🔍 Find() sin filtros retorna todos los documentos de la colección
	// bson.M{} es un filtro vacío (equivale a {} en DB shell)
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return domain.PaginatedResponse{}, err
	}
	defer cur.Close(ctx) // ⚠️ IMPORTANTE: Siempre cerrar el cursor para liberar recursos

	// 📦 Decodificar resultados en slice de DAO (modelo DB)
	// Usamos el modelo DAO porque maneja ObjectID y tags BSON
	var daoItems []dao.Item
	if err := cur.All(ctx, &daoItems); err != nil {
		return domain.PaginatedResponse{}, err
	}

	// 🔄 Convertir de DAO a Domain (para la capa de negocio)
	// Separamos los modelos: DAO para DB, Domain para lógica de negocio
	domainItems := make([]domain.Item, len(daoItems))
	for i, daoItem := range daoItems {
		domainItems[i] = daoItem.ToDomain() // Función definida en dao/Item.go
	}

	return domain.PaginatedResponse{
		Page:    1,
		Count:   len(domainItems),
		Results: domainItems,
	}, nil
}

// Create inserta un nuevo item en DB
// Consigna 1: Validar name y price >= 0, agregar timestamps
func (r *MongoItemsRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	itemDAO := dao.FromDomain(item) // Convertir a DAO para manejar ObjectID y BSON
	itemDAO.ID = primitive.NewObjectID()
	itemDAO.CreatedAt = time.Now().UTC()
	itemDAO.UpdatedAt = time.Now().UTC()

	// Insertar en DB
	res, err := r.col.InsertOne(ctx, itemDAO)
	if err != nil {
		return domain.Item{}, err
	}

	// Obtener el ID generado por DB y convertir a string
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		itemDAO.ID = oid
	} else {
		return domain.Item{}, errors.New("failed to convert inserted ID to ObjectID")
	}

	return itemDAO.ToDomain(), nil // Convertir de vuelta a Domain para retornar
}

// GetByID busca un item por su ID
// Consigna 2: Validar que el ID sea un ObjectID válido
func (r *MongoItemsRepository) GetByID(ctx context.Context, id string) (domain.Item, error) {
	// Validar que el ID tenga formato ObjectID válido
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Item{}, errors.New("invalid ObjectID format")
	}

	// Buscar el documento por _id
	var daoItem dao.Item
	filter := bson.M{"_id": objectID}
	err = r.col.FindOne(ctx, filter).Decode(&daoItem)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Item{}, errors.New("item not found")
		}
		return domain.Item{}, err
	}

	return daoItem.ToDomain(), nil
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
