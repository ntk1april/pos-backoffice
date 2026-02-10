import apiClient from "./client";
import { LoginRequest, LoginResponse, ApiResponse } from "../types";

export const authApi = {
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<ApiResponse<LoginResponse>>(
      "/auth/login",
      credentials,
    );
    return response.data.data!;
  },
};
