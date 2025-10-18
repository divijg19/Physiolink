// src/screens/ProfileScreen.js
import React, { useContext, useCallback } from "react";
import {
	View,
	Text,
	StyleSheet,
	ActivityIndicator,
	Alert,
	ScrollView,
} from "react-native";
import { useFocusEffect } from "@react-navigation/native"; // 1. Import useFocusEffect
import StyledButton from "../components/StyledButton";
import { AuthContext } from "../context/AuthContext";
import { COLORS } from "../theme";
import apiClient from "../api/client";

const ProfileScreen = ({ navigation }) => {
	// 2. Add navigation prop
	const { signOut, userRole } = useContext(AuthContext);
	const [profile, setProfile] = React.useState(null);
	const [isLoading, setIsLoading] = React.useState(true);

	// 3. This is a smart hook that re-runs the fetch logic every time the screen comes into focus
	useFocusEffect(
		useCallback(() => {
			const fetchProfile = async () => {
				setIsLoading(true);
				try {
					const response = await apiClient.get("/profile/me");
					setProfile(response.data);
				} catch (error) {
					// This is okay, it just means the user has no profile yet
					setProfile(null);
				} finally {
					setIsLoading(false);
				}
			};
			fetchProfile();
		}, []),
	);

	if (isLoading) {
		return (
			<View style={styles.container}>
				<ActivityIndicator size="large" color={COLORS.primary} />
			</View>
		);
	}

	// 4. Update the UI to navigate to the EditProfileScreen
	if (profile) {
		// User HAS a profile
		return (
			<ScrollView contentContainerStyle={styles.container}>
				<Text style={styles.title}>Welcome, {profile.firstName}</Text>
				<View style={styles.infoBox}>
					<Text style={styles.infoLabel}>Email:</Text>
					<Text style={styles.infoText}>{profile.user?.email}</Text>
				</View>

				{/* Role-aware display: show PT fields for 'pt', patient fields otherwise */}
				{userRole === "pt" ? (
					<>
						<View style={styles.infoBox}>
							<Text style={styles.infoLabel}>Specialty:</Text>
							<Text style={styles.infoText}>
								{profile.specialty || "Not specified"}
							</Text>
						</View>
						<View style={styles.infoBox}>
							<Text style={styles.infoLabel}>Bio:</Text>
							<Text style={styles.infoText}>
								{profile.bio || "Not specified"}
							</Text>
						</View>
					</>
				) : (
					<>
						<View style={styles.infoBox}>
							<Text style={styles.infoLabel}>Age:</Text>
							<Text style={styles.infoText}>{profile.age || 'Not specified'}</Text>
						</View>
						<View style={styles.infoBox}>
							<Text style={styles.infoLabel}>Gender:</Text>
							<Text style={styles.infoText}>{profile.gender || 'Not specified'}</Text>
						</View>
						<View style={styles.infoBox}>
							<Text style={styles.infoLabel}>Condition:</Text>
							<Text style={styles.infoText}>
								{profile.condition || "Not specified"}
							</Text>
						</View>
						<View style={styles.infoBox}>
							<Text style={styles.infoLabel}>Goals:</Text>
							<Text style={styles.infoText}>
								{profile.goals || "Not specified"}
							</Text>
						</View>
					</>
				)}

				<StyledButton
					title="Edit Profile"
					onPress={() => navigation.navigate("EditProfile", { profile })}
				/>
				<StyledButton title="Sign Out" onPress={signOut} />
			</ScrollView>
		);
	} else {
		// User does NOT have a profile
		return (
			<View style={styles.container}>
				<Text style={styles.title}>Welcome!</Text>
				<Text style={styles.subtitle}>Let's create your profile.</Text>
				<StyledButton
					title="Create Profile"
					onPress={() => navigation.navigate("EditProfile", { profile: null })}
				/>
				<StyledButton title="Sign Out" onPress={signOut} />
			</View>
		);
	}
};

// --- (The styles)
const styles = StyleSheet.create({
	container: {
		flex: 1,
		justifyContent: "center",
		alignItems: "center",
		padding: 20,
		backgroundColor: COLORS.background,
	},
	title: {
		fontFamily: "Poppins_600SemiBold",
		fontSize: 26,
		color: COLORS.textDark,
		marginBottom: 20,
		textAlign: "center",
	},
	subtitle: {
		fontFamily: "Poppins_400Regular",
		fontSize: 16,
		color: COLORS.gray,
		marginBottom: 20,
	},
	infoBox: {
		width: "100%",
		backgroundColor: COLORS.lightGray,
		padding: 15,
		borderRadius: 10,
		marginBottom: 10,
	},
	infoLabel: {
		fontFamily: "Poppins_600SemiBold",
		fontSize: 14,
		color: COLORS.gray,
	},
	infoText: {
		fontFamily: "Poppins_400Regular",
		fontSize: 16,
		color: COLORS.textDark,
	},
});

export default ProfileScreen;
