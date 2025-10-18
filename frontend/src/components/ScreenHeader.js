// src/components/ScreenHeader.js
import React from "react";
import { View, Text, StyleSheet } from "react-native";
import { COLORS } from "../theme";

const ScreenHeader = ({ title, subtitle }) => (
	<View style={styles.headerContainer}>
		<Text style={styles.header}>{title}</Text>
		{subtitle && <Text style={styles.subtitle}>{subtitle}</Text>}
	</View>
);

const styles = StyleSheet.create({
	headerContainer: {
		paddingHorizontal: 20,
		paddingTop: 20,
		paddingBottom: 10,
		backgroundColor: COLORS.background,
	},
	header: {
		fontFamily: "Poppins_700Bold",
		fontSize: 28,
		color: COLORS.textDark,
	},
	subtitle: {
		fontFamily: "Poppins_400Regular",
		fontSize: 16,
		color: COLORS.gray,
		marginTop: 4,
	},
});

export default ScreenHeader;
