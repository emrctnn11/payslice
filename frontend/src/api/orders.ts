const API_BASE_URL = "http://localhost:8080";

interface CreateOrderRequest {
  product_id: number;
  user_id: string; //we dont have an user token rn.
}

interface OrderResponse {
  id: string;
  status: string;
  total_amount_cents: number;
}

export const createOrder = async (
  productId: number
): Promise<OrderResponse> => {
  // For this MVP, we fake a userID.
  const payload: CreateOrderRequest = {
    product_id: productId,
    user_id: "demo-user-1",
  };

  const response = await fetch(`${API_BASE_URL}/orders`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    throw new Error(`Order Failed: ${response.statusText}`);
  }

  return await response.json();
};
