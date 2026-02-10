import axios from "axios";

const API_URL = "http://localhost:8080/api";

// Get auth token from localStorage
const getAuthHeader = () => {
  const token = localStorage.getItem("token");
  return token ? { Authorization: `Bearer ${token}` } : {};
};

export interface Store {
  id: number;
  code: string;
  name: string;
  address: string;
  phone: string;
  status: string;
  created_at: string;
  updated_at: string;
  created_by: number;
  updated_by: number;
}

export interface CreateStoreRequest {
  code: string;
  name: string;
  address?: string;
  phone?: string;
  status?: string;
}

export interface UpdateStoreRequest {
  code: string;
  name: string;
  address?: string;
  phone?: string;
  status?: string;
}

export const storeApi = {
  // Get all stores
  getStores: async (): Promise<Store[]> => {
    const response = await axios.get(`${API_URL}/stores`, {
      headers: getAuthHeader(),
    });
    return response.data.data;
  },

  // Get store by ID
  getStore: async (id: number): Promise<Store> => {
    const response = await axios.get(`${API_URL}/stores/${id}`, {
      headers: getAuthHeader(),
    });
    return response.data.data;
  },

  // Create store
  createStore: async (data: CreateStoreRequest): Promise<Store> => {
    const response = await axios.post(`${API_URL}/stores`, data, {
      headers: getAuthHeader(),
    });
    return response.data.data;
  },

  // Update store
  updateStore: async (id: number, data: UpdateStoreRequest): Promise<Store> => {
    const response = await axios.put(`${API_URL}/stores/${id}`, data, {
      headers: getAuthHeader(),
    });
    return response.data.data;
  },

  // Delete store
  deleteStore: async (id: number): Promise<void> => {
    await axios.delete(`${API_URL}/stores/${id}`, {
      headers: getAuthHeader(),
    });
  },
};
