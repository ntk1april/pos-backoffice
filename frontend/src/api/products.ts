import apiClient from "./client";
import {
  Product,
  ProductListResponse,
  CreateProductRequest,
  UpdateProductRequest,
  StockAdjustmentRequest,
  StockLogResponse,
  ApiResponse,
} from "../types";

export const productApi = {
  getProducts: async (
    page: number = 1,
    pageSize: number = 10,
    search: string = "",
    status: string = "ACTIVE",
  ): Promise<ProductListResponse> => {
    const response = await apiClient.get<ApiResponse<ProductListResponse>>(
      "/products",
      {
        params: { page, page_size: pageSize, search, status },
      },
    );
    return response.data.data!;
  },

  getProduct: async (id: number): Promise<Product> => {
    const response = await apiClient.get<ApiResponse<Product>>(
      `/products/${id}`,
    );
    return response.data.data!;
  },

  createProduct: async (product: CreateProductRequest): Promise<Product> => {
    const response = await apiClient.post<ApiResponse<Product>>(
      "/products",
      product,
    );
    return response.data.data!;
  },

  updateProduct: async (
    id: number,
    product: UpdateProductRequest,
  ): Promise<Product> => {
    const response = await apiClient.put<ApiResponse<Product>>(
      `/products/${id}`,
      product,
    );
    return response.data.data!;
  },

  deleteProduct: async (id: number): Promise<void> => {
    await apiClient.delete(`/products/${id}`);
  },

  increaseStock: async (request: StockAdjustmentRequest): Promise<void> => {
    await apiClient.post("/stock/increase", request);
  },

  decreaseStock: async (request: StockAdjustmentRequest): Promise<void> => {
    await apiClient.post("/stock/decrease", request);
  },

  getStockLogs: async (
    productId: number,
    page: number = 1,
    pageSize: number = 20,
  ): Promise<StockLogResponse> => {
    const response = await apiClient.get<ApiResponse<StockLogResponse>>(
      `/stock/logs/${productId}`,
      {
        params: { page, page_size: pageSize },
      },
    );
    return response.data.data!;
  },
};
