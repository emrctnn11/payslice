import type { Product } from "../types";

const API_BASE_URL = "http://localhost:8080";

export const getProducts = async (): Promise<Product[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/products`);
    if (!response.ok) {
      throw new Error(`Error fetching products: ${response.statusText}`);
    }

    return await response.json();
  } catch (error) {
        console.error('Api Error:', error);
        throw error;
  }
};
