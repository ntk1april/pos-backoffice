export interface User {
  id: number;
  username: string;
  full_name: string;
  role: "ADMIN" | "STAFF";
  status: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface Product {
  id: number;
  sku: string;
  name: string;
  description: string;
  price: number;
  cost: number;
  stock: number;
  status: "ACTIVE" | "INACTIVE";
  created_at: string;
  updated_at: string;
  created_by: number;
  updated_by: number;
}

export interface CreateProductRequest {
  sku: string;
  name: string;
  description: string;
  price: number;
  cost: number;
  stock: number;
}

export interface UpdateProductRequest {
  name: string;
  description: string;
  price: number;
  cost: number;
}

export interface ProductListResponse {
  products: Product[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface StockAdjustmentRequest {
  product_id: number;
  quantity: number;
  notes: string;
}

export interface StockLog {
  id: number;
  product_id: number;
  transaction_type: "INCREASE" | "DECREASE" | "ADJUSTMENT";
  quantity: number;
  stock_before: number;
  stock_after: number;
  notes: string;
  created_by: number;
  created_at: string;
}

export interface StockLogResponse {
  logs: StockLog[];
  total: number;
  page: number;
  page_size: number;
}

export interface ApiResponse<T> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}
