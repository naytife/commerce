import { dev } from '$app/environment';

export interface GraphQLError {
  message: string;
  path?: string[];
  code?: string;
}

export interface GraphQLResponse<T = any> {
  data?: T;
  errors?: GraphQLError[];
}

export class ApiClient {
  private restUrl: string;
  private debug: boolean;

  constructor() {
    // Use REST API endpoint for REST calls
    this.restUrl = import.meta.env.VITE_API_URL || 'http://localhost:8002/api/v1';
    this.debug = import.meta.env.VITE_DEBUG === 'true' || dev;
  }

  async restGet<T = any>(endpoint: string): Promise<T> {
    try {
      const url = `${this.restUrl}${endpoint}`;
      
      if (this.debug) {
        console.log('REST GET:', url);
      }

      const response = await fetch(url);

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const result = await response.json();

      if (this.debug) {
        console.log('REST Response:', result);
      }

      // Handle backend API response format
      if (result.data) {
        return result.data;
      }
      
      return result;
    } catch (error) {
      console.error('REST API Error:', error);
      throw error;
    }
  }

  async restPost<T = any>(endpoint: string, data: any): Promise<T> {
    try {
      const url = `${this.restUrl}${endpoint}`;
      
      if (this.debug) {
        console.log('REST POST:', url, data);
      }

      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const result = await response.json();

      if (this.debug) {
        console.log('REST Response:', result);
      }

      // Handle backend API response format
      if (result.data) {
        return result.data;
      }
      
      return result;
    } catch (error) {
      console.error('REST API Error:', error);
      throw error;
    }
  }
}

// Singleton instance
export const apiClient = new ApiClient();
