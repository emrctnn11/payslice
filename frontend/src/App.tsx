import { useEffect, useState } from "react";
import "./App.css";
import type { Product } from "./types";
import { getProducts } from "./api/products";
import { PorductCard } from "./components/ProductCard";
import { createOrder } from "./api/orders";

function App() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [buyingId, setBuyingId] = useState<number | null>(null);

  useEffect(() => {
    const loadData = async () => {
      try {
        const data = await getProducts();
        setProducts(data);
      } catch {
        setError("Failed");
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  const handleBuy = async (product: Product) => {
    console.log("Initiating purchase for: ", product.name);
    try {
      setBuyingId(product.id); //loading spinner on button

      const order = await createOrder(product.id);

      alert(
        `Order Placed Successfully!\nOrder ID: ${order.id}\nStatus: ${order.status}`
      );

      setProducts((currentProducts) =>
        currentProducts.map(
          (p) =>
            p.id === product.id ? { ...p, inventory: p.inventory - 1 } : p //
        )
      );
    } catch (err) {
      console.error(err);
      alert("failed to place order.");
    } finally {
      setBuyingId(null);
    }
  };
  if (loading)
    return (
      <div className="text-center p-10 text-gray-500">Loading Store...</div>
    );
  if (error)
    return (
      <div className="text-center p-10 text-red-500 font-bold">{error}</div>
    );

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <header className="mb-12 text-center max-w-4xl mx-auto">
        <h1 className="text-4xl font-extrabold text-gray-900 tracking-tight mb-2">
          PaySlice Store 🍕
        </h1>
        <p className="text-gray-500">Buy now, pay later. Zero friction.</p>
      </header>

      <div className="max-w-6xl mx-auto grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {products.map((product) => (
          <PorductCard
            key={product.id}
            product={product}
            onBuy={handleBuy}
            isBuying={buyingId === product.id}
          />
        ))}
      </div>
    </div>
  );
}

export default App;
