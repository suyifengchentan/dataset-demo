db = db.getSiblingDB("visual_lab");

db.orders.updateOne(
  { _id: 1 },
  {
    $set: {
      customer: { name: "Alice", level: "vip" },
      status: "paid",
      tags: ["vip", "web"],
      items: [
        { sku: "A-1", category: "sensor", qty: 2 },
        { sku: "B-1", category: "cable", qty: 4 }
      ],
      shipping: { city: "Shanghai", method: "air" }
    }
  },
  { upsert: true }
);

db.orders.updateOne(
  { _id: 2 },
  {
    $set: {
      customer: { name: "Bob", level: "standard" },
      status: "pending",
      tags: ["mobile"],
      items: [
        { sku: "C-2", category: "module", qty: 1 }
      ],
      shipping: { city: "Beijing", method: "ground" }
    }
  },
  { upsert: true }
);
