// src/screens/RegisterScreen.js
import React, { useState, useContext } from "react";
import { View, Text, StyleSheet, Alert, TouchableOpacity } from "react-native";
// ====================================================================
// THE FIX: Import SafeAreaView from the correct library
import { SafeAreaView } from "react-native-safe-area-context";
// ====================================================================
import StyledInput from "../components/StyledInput";
import StyledButton from "../components/StyledButton";
import { COLORS } from "../theme";
import apiClient from "../api/client";
import { AuthContext } from "../context/AuthContext";

const RegisterScreen = ({ navigation }) => {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const [confirmPassword, setConfirmPassword] = useState("");
	const [isLoading, setIsLoading] = useState(false);
	const { signIn } = useContext(AuthContext);

	const handleRegister = async () => {
		if (!email || !password || !confirmPassword) {
			return Alert.alert("Missing Fields", "Please fill out all fields.");
		}
		if (password !== confirmPassword) {
			return Alert.alert("Password Mismatch", "The passwords do not match.");
		}

		setIsLoading(true);
		try {
			const response = await apiClient.post("/auth/register", {
				email,
				password,
				role: "patient",
			});
			const { token, profile } = response.data;
			// Save session
			signIn(token);
			// If profile exists but is empty, navigate to EditProfile to complete it
			if (profile && (!profile.firstName || !profile.lastName)) {
				navigation.navigate('EditProfile', { profile });
			}
		} catch (error) {
			const errorMsg =
				error.response?.data?.msg || "An unexpected error occurred.";
			Alert.alert("Registration Failed", errorMsg);
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<SafeAreaView style={styles.container}>
			<Text style={styles.title}>Create Account</Text>
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
			<StyledInput
				placeholder="Confirm Password"
				value={confirmPassword}
				onChangeText={setConfirmPassword}
				isPassword
			/>
			<StyledButton
				title="Register"
				onPress={handleRegister}
				isLoading={isLoading}
			/>
			<TouchableOpacity onPress={() => navigation.goBack()}>
				<Text style={styles.linkText}>Already have an account? Log In</Text>
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
	linkText: {
		color: COLORS.primary,
		textAlign: "center",
		marginTop: 20,
		fontFamily: "Poppins_600SemiBold",
	},
});

export default RegisterScreen;
