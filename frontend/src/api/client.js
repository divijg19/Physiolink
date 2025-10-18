// src/api/client.js
import axios from "axios";
import * as SecureStore from "expo-secure-store";

// This IP address is correct from your file.
const API_BASE_URL = "http://192.168.1.15:3001/api";

const apiClient = axios.create({
	baseURL: API_BASE_URL,
});

// ====================================================================
// CRITICAL FIX: Add a request interceptor to attach the token
// ====================================================================
apiClient.interceptors.request.use(
	async (config) => {
		// Get the token from secure storage
		const token = await SecureStore.getItemAsync("userToken");
		if (token) {
			// If a token exists, add it to the 'x-auth-token' header
			config.headers["x-auth-token"] = token;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	},
);
// ====================================================================

export default apiClient;
