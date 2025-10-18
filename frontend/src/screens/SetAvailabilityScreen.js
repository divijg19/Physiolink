// src/screens/SetAvailabilityScreen.js
import React, { useState } from "react";
import { View, Text, StyleSheet, Alert } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import StyledButton from "../components/StyledButton";
import { COLORS } from "../theme";
import apiClient from "../api/client";
import ScreenHeader from "../components/ScreenHeader";

const SetAvailabilityScreen = () => {
	const [isLoading, setIsLoading] = useState(false);

	const handleCreateDemoSlots = async () => {
		setIsLoading(true);
		try {
			const today = new Date();
			const tomorrow = new Date(today);
			tomorrow.setDate(tomorrow.getDate() + 1);

			const slots = [
				{
					startTime: new Date(tomorrow.setHours(9, 0, 0, 0)),
					endTime: new Date(tomorrow.setHours(10, 0, 0, 0)),
				},
				{
					startTime: new Date(tomorrow.setHours(10, 0, 0, 0)),
					endTime: new Date(tomorrow.setHours(11, 0, 0, 0)),
				},
				{
					startTime: new Date(tomorrow.setHours(11, 0, 0, 0)),
					endTime: new Date(tomorrow.setHours(12, 0, 0, 0)),
				},
			];

			await apiClient.post("/appointments/availability", { slots });
			Alert.alert(
				"Success",
				"Your available slots have been created for tomorrow.",
			);
		} catch (error) {
			const errorMsg = error.response?.data?.msg || "Could not create slots.";
			Alert.alert("Error", errorMsg);
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<SafeAreaView style={styles.container}>
			<ScreenHeader title="Set Your Availability" />
			<View style={styles.content}>
				<Text style={styles.description}>
					This is a demo feature to add availability for your clients to book.
				</Text>
				<StyledButton
					title="Create Demo Slots for Tomorrow"
					onPress={handleCreateDemoSlots}
					isLoading={isLoading}
				/>
			</View>
		</SafeAreaView>
	);
};

const styles = StyleSheet.create({
	container: { flex: 1, backgroundColor: COLORS.background },
	content: { flex: 1, justifyContent: "center", padding: 20 },
	description: {
		fontFamily: "Poppins_400Regular",
		fontSize: 16,
		textAlign: "center",
		color: COLORS.textDark,
		marginBottom: 30,
	},
});

export default SetAvailabilityScreen;
