// src/screens/ProfileScreen.js
import React, { useContext, useCallback } from "react";
import {
	View,
	Text,
	StyleSheet,
	ActivityIndicator,
	ScrollView,
} from "react-native";
import Avatar from '../components/Avatar';
import { useFocusEffect } from "@react-navigation/native"; // 1. Import useFocusEffect
import StyledButton from "../components/StyledButton";
import { AuthContext } from "../context/AuthContext";
import { COLORS, SPACING, FONT, SHADOW } from "../theme";
import Snackbar from '../components/Snackbar';
import apiClient from "../api/client";

const ProfileScreen = ({ navigation }) => {
	// 2. Add navigation prop
	const { signOut, userRole } = useContext(AuthContext);
	const [profile, setProfile] = React.useState(null);
	const [isLoading, setIsLoading] = React.useState(true);
	const [snack, setSnack] = React.useState({ visible: false, message: '' });

	// 3. This is a smart hook that re-runs the fetch logic every time the screen comes into focus
	useFocusEffect(
		useCallback(() => {
			const fetchProfile = async () => {
				setIsLoading(true);
				try {
					const response = await apiClient.get("/profile/me");
					setProfile(response.data);
				} catch (error) {
					// show a subtle message but allow flow to continue (user may not have a profile yet)
					setProfile(null);
					setSnack({ visible: true, message: 'Could not load profile.' });
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
				<Avatar uri={profile.profileImageUrl} name={`${profile.firstName} ${profile.lastName}`} size={96} />
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
							<Text style={styles.infoLabel}>Location:</Text>
							<Text style={styles.infoText}>{profile.location || 'Not specified'}</Text>
						</View>
						{profile.credentials ? (
							<View style={styles.infoBox}>
								<Text style={styles.infoLabel}>Credentials:</Text>
								<Text style={styles.infoText}>{profile.credentials}</Text>
							</View>
						) : null}
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

				<View style={styles.buttonsRow}>
					<StyledButton
						style={styles.singleButton}
						title="Edit Profile"
						onPress={() => navigation.navigate("EditProfile", { profile })}
					/>
					<StyledButton style={styles.singleButton} title="Sign Out" onPress={signOut} />
				</View>
			</ScrollView>
		);
	} else {
		// User does NOT have a profile
		return (
			<View style={styles.container}>
				<Text style={styles.title}>Welcome!</Text>
				<Text style={styles.subtitle}>Let's create your profile.</Text>
				<View style={styles.buttonsRow}>
					<StyledButton style={styles.singleButton}
						title="Create Profile"
						onPress={() => navigation.navigate("EditProfile", { profile: null })}
					/>
					<StyledButton style={styles.singleButton} title="Sign Out" onPress={signOut} />
				</View>
				<Snackbar visible={snack.visible} message={snack.message} />
			</View>
		);
	}
};

// --- (The styles)
const styles = StyleSheet.create({
	container: {
		flex: 1,
		padding: SPACING.md,
		backgroundColor: COLORS.background,
		alignItems: 'center',
	},
	title: {
		fontFamily: "Poppins_600SemiBold",
		fontSize: FONT.largeTitle,
		color: COLORS.textDark,
		marginBottom: SPACING.md,
		textAlign: "center",
	},
	subtitle: {
		fontFamily: "Poppins_400Regular",
		fontSize: FONT.body,
		color: COLORS.gray,
		marginBottom: SPACING.md,
	},
	infoBox: {
		width: "100%",
		backgroundColor: COLORS.lightGray,
		padding: SPACING.md,
		borderRadius: 10,
		marginBottom: SPACING.sm,
		...SHADOW,
	},
	infoLabel: {
		fontFamily: "Poppins_600SemiBold",
		fontSize: FONT.small,
		color: COLORS.gray,
	},
	infoText: {
		fontFamily: "Poppins_400Regular",
		fontSize: FONT.body,
		color: COLORS.textDark,
	},
	buttonsRow: { width: '100%', marginTop: SPACING.md },
	singleButton: { marginBottom: SPACING.sm },
});

export default ProfileScreen;
