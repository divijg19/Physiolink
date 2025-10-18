// src/components/StyledButton.js
import React from "react";
import {
	TouchableOpacity,
	Text,
	StyleSheet,
	ActivityIndicator,
} from "react-native";
import { COLORS } from "../theme";

// Add the isLoading prop
const StyledButton = ({ title, onPress, isLoading = false }) => {
	return (
		<TouchableOpacity
			style={styles.button}
			onPress={onPress}
			disabled={isLoading}
		>
			{isLoading ? (
				<ActivityIndicator color={COLORS.textLight} />
			) : (
				<Text style={styles.buttonText}>{title}</Text>
			)}
		</TouchableOpacity>
	);
};

const styles = StyleSheet.create({
	// ... same styles as before
	button: {
		backgroundColor: COLORS.accent,
		padding: 18,
		borderRadius: 12,
		alignItems: "center",
		marginVertical: 10,
		justifyContent: "center",
		minHeight: 60,
	},
	buttonText: {
		color: COLORS.textLight,
		fontSize: 18,
		fontFamily: "Poppins_600SemiBold",
	},
});

export default StyledButton;
