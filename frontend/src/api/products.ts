import apiClient from "./client";
import {
  Product,
  ProductListResponse,
  CreateProductRequest,
  UpdateProductRequest,
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
};
