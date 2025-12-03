// src/components/StyledButton.js
import React from "react";
import {
	TouchableOpacity,
	Text,
	StyleSheet,
	ActivityIndicator,
	Platform,
} from "react-native";
import { COLORS, SPACING, FONT, SHADOW } from "../theme";

// Add the isLoading prop, disabled, and size
const StyledButton = ({ title, onPress, isLoading = false, variant = 'primary', style, disabled = false, size = 'md' }) => {
	const isGhost = variant === 'ghost';
	const colorMap = {
		primary: COLORS.primary,
		secondary: COLORS.secondary,
		accent: COLORS.accent,
	};
	const bgColor = isGhost ? 'transparent' : (colorMap[variant] || COLORS.primary);

	const sizeStyle = size === 'sm' ? styles.small : null;
	const disabledStyle = disabled ? styles.disabled : null;

	const buttonStyles = [styles.button, { backgroundColor: bgColor }, isGhost && styles.ghost, sizeStyle, disabledStyle, style];

	const indicatorColor = isGhost ? COLORS.primary : COLORS.textLight;

	return (
		<TouchableOpacity
			style={buttonStyles}
			onPress={onPress}
			disabled={isLoading || disabled}
			activeOpacity={0.8}
		>
			{isLoading ? (
				<ActivityIndicator color={indicatorColor} />
			) : (
				<Text style={[styles.buttonText, isGhost ? styles.ghostText : null]}>{title}</Text>
			)}
		</TouchableOpacity>
	);
};

const styles = StyleSheet.create({
	button: {
		backgroundColor: COLORS.accent,
		paddingVertical: SPACING.sm,
		paddingHorizontal: SPACING.md,
		borderRadius: 12,
		alignItems: "center",
		marginVertical: SPACING.sm,
		justifyContent: "center",
		minHeight: 48,
		flexDirection: 'row',
		justifyContent: 'center',
		...(Platform.OS === 'android' ? { elevation: SHADOW.elevation } : { shadowColor: SHADOW.shadowColor, shadowOffset: SHADOW.shadowOffset, shadowOpacity: SHADOW.shadowOpacity, shadowRadius: SHADOW.shadowRadius }),
	},
	small: {
		paddingVertical: SPACING.xs,
		minHeight: 40,
	},
	disabled: {
		opacity: 0.7,
	},
	ghost: {
		backgroundColor: 'transparent',
		borderWidth: 1,
		borderColor: COLORS.primary,
	},
	buttonText: {
		color: COLORS.textLight,
		fontSize: FONT.body,
		fontFamily: "Poppins_600SemiBold",
	},
	ghostText: {
		color: COLORS.primary,
	},
});

export default StyledButton;
