// Seed some data if the collection is empty
db = db.getSiblingDB(process.env.MONGO_INITDB_DATABASE || "demo");
const col = db.getCollection("items");
if (col.countDocuments() === 0) {
  col.insertMany([
    { name: "Notebook", price: 12.5, createdAt: new Date() },
    { name: "Coffee", price: 3.25, createdAt: new Date() },
    { name: "Keyboard", price: 22.0, createdAt: new Date() }
  ]);
  print("Seeded initial items");
} else {
  print("Items already present, skipping seed");
}
