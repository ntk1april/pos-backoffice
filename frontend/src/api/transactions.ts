import axios from "axios";

const API_URL = "http://localhost:8080/api";

// Get auth token from localStorage
const getAuthHeader = () => {
  const token = localStorage.getItem("token");
  return token ? { Authorization: `Bearer ${token}` } : {};
};

export interface Transaction {
  id: number;
  transaction_type: "INCREASE" | "DECREASE";
  product_id: number;
  product_name?: string;
  store_id?: number | null;
  store_name?: string;
  quantity: number;
  unit_price: number;
  total_amount: number;
  notes: string;
  transaction_date: string;
  created_by: number;
  created_by_name?: string;
}

export interface CreateTransactionRequest {
  transaction_type: "INCREASE" | "DECREASE";
  product_id: number;
  store_id?: number | null;
  quantity: number;
  unit_price: number;
  notes?: string;
}

export interface TransactionListResponse {
  transactions: Transaction[];
  total: number;
  page: number;
  limit: number;
}

export const transactionApi = {
  // Create transaction (INCREASE or DECREASE)
  createTransaction: async (
    data: CreateTransactionRequest,
  ): Promise<Transaction> => {
    const response = await axios.post(`${API_URL}/transactions`, data, {
      headers: getAuthHeader(),
    });
    return response.data.data;
  },

  // Get all transactions with pagination
  getTransactions: async (
    page: number = 1,
    limit: number = 20,
  ): Promise<TransactionListResponse> => {
    const response = await axios.get(`${API_URL}/transactions`, {
      params: { page, limit },
      headers: getAuthHeader(),
    });
    return response.data.data;
  },

  // Get transactions by product
  getTransactionsByProduct: async (
    productId: number,
    limit: number = 50,
  ): Promise<Transaction[]> => {
    const response = await axios.get(
      `${API_URL}/transactions/product/${productId}`,
      {
        params: { limit },
        headers: getAuthHeader(),
      },
    );
    return response.data.data;
  },

  // Get transactions by store
  getTransactionsByStore: async (
    storeId: number,
    limit: number = 50,
  ): Promise<Transaction[]> => {
    const response = await axios.get(
      `${API_URL}/transactions/store/${storeId}`,
      {
        params: { limit },
        headers: getAuthHeader(),
      },
    );
    return response.data.data;
  },
};
