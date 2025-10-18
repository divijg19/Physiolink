// src/screens/LoginScreen.js
import React, { useContext, useState } from "react";
import { Alert, StyleSheet, Text, TouchableOpacity } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import apiClient from "../api/client";
import StyledButton from "../components/StyledButton";
import StyledInput from "../components/StyledInput";
import { AuthContext } from "../context/AuthContext";
import { COLORS } from "../theme";

// The navigation prop is passed down from the Stack Navigator
const LoginScreen = ({ navigation }) => {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	// 1. Add loading state
	const [isLoading, setIsLoading] = useState(false);
	const { signIn } = useContext(AuthContext);

	const handleLogin = async () => {
		if (!email || !password) return;
		// 2. Set loading to true when the process starts
		setIsLoading(true);
		try {
			const response = await apiClient.post("/auth/login", { email, password });
			signIn(response.data.token);
		} catch (_error) {
			Alert.alert("Login Failed", "Invalid credentials. Please try again.");
		} finally {
			// 3. Set loading to false when the process finishes
			setIsLoading(false);
		}
	};

	return (
		<SafeAreaView style={styles.container}>
			<Text style={styles.title}>Welcome Back</Text>
			<StyledInput
				placeholder="Email"
				value={email}
				onChangeText={setEmail}
				keyboardType="email-address"
				autoCapitalize="none"
			/>
			<StyledInput
				placeholder="Password"
				value={password}
				onChangeText={setPassword}
				isPassword
			/>
			{/* 4. Pass the isLoading state to the button */}
			<StyledButton
				title="Log In"
				onPress={handleLogin}
				isLoading={isLoading}
			/>

			{/* 5. Add a link to navigate to the Register screen */}
			<TouchableOpacity onPress={() => navigation.navigate("Register")}>
				<Text style={styles.linkText}>Don't have an account? Sign Up</Text>
			</TouchableOpacity>
		</SafeAreaView>
	);
};

const styles = StyleSheet.create({
	container: {
		flex: 1,
		justifyContent: "center",
		padding: 20,
		backgroundColor: COLORS.background,
	},
	title: {
		fontSize: 32,
		fontFamily: "Poppins_700Bold",
		color: COLORS.textDark,
		textAlign: "center",
		marginBottom: 20,
	},
	// 6. Add style for the new link
	linkText: {
		color: COLORS.primary,
		textAlign: "center",
		marginTop: 20,
		fontFamily: "Poppins_600SemiBold",
		fontSize: 16,
	},
});

export default LoginScreen;
