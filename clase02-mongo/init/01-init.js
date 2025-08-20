const appDb = db.getSiblingDB('app');
appDb.createUser({ user: 'appuser', pwd: 'apppass', roles: [{ role: 'readWrite', db: 'app' }] });
appDb.createCollection('items');
appDb.items.createIndex({ name: 1 }, { unique: true });

// 📚 Datos de ejemplo para los estudiantes
const now = new Date();
appDb.items.insertMany([
  {
    name: "Laptop Gaming",
    price: 1299.99,
    created_at: new Date(now.getTime() - 48 * 60 * 60 * 1000), // 48h atrás
    updated_at: new Date(now.getTime() - 24 * 60 * 60 * 1000)  // 24h atrás
  },
  {
    name: "Mouse Inalámbrico",
    price: 29.99,
    created_at: new Date(now.getTime() - 24 * 60 * 60 * 1000), // 24h atrás
    updated_at: new Date(now.getTime() - 12 * 60 * 60 * 1000)  // 12h atrás
  },
  {
    name: "Teclado Mecánico",
    price: 89.50,
    created_at: new Date(now.getTime() - 12 * 60 * 60 * 1000), // 12h atrás
    updated_at: new Date(now.getTime() - 6 * 60 * 60 * 1000)   // 6h atrás
  },
  {
    name: "Monitor 4K",
    price: 299.00,
    created_at: new Date(now.getTime() - 6 * 60 * 60 * 1000),  // 6h atrás
    updated_at: new Date(now.getTime() - 1 * 60 * 60 * 1000)   // 1h atrás
  }
]);

print("✅ Datos de ejemplo insertados correctamente");
