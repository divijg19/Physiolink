// src/context/AuthContext.js
import React, { createContext, useState, useEffect } from "react";
import * as SecureStore from "expo-secure-store";
import { jwtDecode } from "jwt-decode"; // 1. Import the decoder

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
	const [userToken, setUserToken] = useState(null);
	const [userRole, setUserRole] = useState(null); // 2. Add state for user role
	const [isLoading, setIsLoading] = useState(true);

	const setSession = async (token) => {
		if (token) {
			await SecureStore.setItemAsync("userToken", token);
			const decodedToken = jwtDecode(token); // 3. Decode the token
			setUserToken(token);
			setUserRole(decodedToken.user.role); // 4. Set the user's role
		} else {
			await SecureStore.deleteItemAsync("userToken");
			setUserToken(null);
			setUserRole(null);
		}
		setIsLoading(false);
	};

	useEffect(() => {
		const bootstrapAsync = async () => {
			const token = await SecureStore.getItemAsync("userToken");
			setSession(token);
		};
		bootstrapAsync();
	}, []);

	const authContext = {
		signIn: (token) => setSession(token),
		signOut: () => setSession(null),
		userToken,
		userRole, // 5. Expose the role to the rest of the app
		isLoading,
	};

	return (
		<AuthContext.Provider value={authContext}>{children}</AuthContext.Provider>
	);
};
