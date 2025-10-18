// src/components/StyledInput.js
import React, { useState } from "react";
import { TextInput, StyleSheet, View, TouchableOpacity } from "react-native";
import { COLORS } from "../theme";
import { Feather } from "@expo/vector-icons";

// The component now accepts an 'isPassword' prop
const StyledInput = ({ isPassword, ...props }) => {
	//It manages its own visibility state
	const [isPasswordVisible, setIsPasswordVisible] = useState(false);

	// If it's not a password field, return the simple input
	if (!isPassword) {
		return (
			<TextInput
				style={styles.input}
				placeholderTextColor={COLORS.gray}
				{...props}
			/>
		);
	}

	// If it IS a password field, return the input with an icon
	return (
		<View style={styles.inputContainer}>
			<TextInput
				style={styles.inputWithIcon}
				placeholderTextColor={COLORS.gray}
				// 3. The secureTextEntry is now controlled by our state
				secureTextEntry={!isPasswordVisible}
				{...props}
			/>
			<TouchableOpacity
				onPress={() => setIsPasswordVisible(!isPasswordVisible)}
			>
				{/* 4. The icon toggles between 'eye' and 'eye-off' */}
				<Feather
					name={isPasswordVisible ? "eye-off" : "eye"}
					size={24}
					color={COLORS.gray}
					style={styles.icon}
				/>
			</TouchableOpacity>
		</View>
	);
};

const styles = StyleSheet.create({
	// Original input style
	input: {
		backgroundColor: COLORS.lightGray,
		paddingVertical: 15,
		paddingHorizontal: 20,
		borderRadius: 12,
		fontSize: 16,
		color: COLORS.textDark,
		marginVertical: 10,
		fontFamily: "Poppins_400Regular",
	},
	// New container for the password input and icon
	inputContainer: {
		flexDirection: "row",
		alignItems: "center",
		backgroundColor: COLORS.lightGray,
		borderRadius: 12,
		marginVertical: 10,
	},
	inputWithIcon: {
		flex: 1, // Take up most of the space
		paddingVertical: 15,
		paddingHorizontal: 20,
		fontSize: 16,
		color: COLORS.textDark,
		fontFamily: "Poppins_400Regular",
	},
	icon: {
		paddingHorizontal: 15,
	},
});

export default StyledInput;
