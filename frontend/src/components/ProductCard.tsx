import type { Product } from "../types";

interface ProductCardProps {
  product: Product;
  onBuy: (product: Product) => void;
  isBuying: boolean;
}

export const PorductCard = ({ product, onBuy, isBuying }: ProductCardProps) => {
  // helper for currency formatter.
  const formatPrice = (cents: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(cents / 100);
  };

  return (
    <div className="bg-white rounded-xl shadow-lg hover:shadow-xl transition-shadow duration-300 overflow-hidden border border-gray-100 flex flex-col">
      {/* image placeholder */}
      <div className="h-48 bg-gray-50 flex items-center justify-center border-b border-gray-100">
        <span className="text-4xl filter grayscale opacity-50">📦</span>
      </div>
      <div className="p-6 flex flex-col flex-grow">
        <h2 className="text-xl font-bold text-gray-800 mb-2 truncate">
          {product.name}
        </h2>

        <div className="flex justify-between items-end mt-auto mb-6">
          <span className="text-2xl font-bold text-indigo-600">
            {formatPrice(product.price_cents)}
          </span>
          <span
            className={`text-sm font-medium ${
              product.inventory > 0 ? "text-green-600" : "text-red-500"
            }`}
          >
            {product.inventory > 0
              ? `${product.inventory} in stock`
              : "Out of Stock"}
          </span>
        </div>

        <button
          onClick={() => onBuy(product)}
          disabled={product.inventory === 0 || isBuying}
          className="w-full bg-indigo-600 text-white py-3 px-4 rounded-lg font-semibold hover:bg-indigo-700 transition-colors disabled:bg-gray-300 disabled:cursor-not-allowed"
        >
          {isBuying ? "Buying..." : "Buy with PaySlice"}
        </button>
      </div>
    </div>
  );
};
