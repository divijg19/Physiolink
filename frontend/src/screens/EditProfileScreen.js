// src/screens/EditProfileScreen.js
import React, { useState, useContext } from "react"; // 1. Import useContext
import { View, Text, StyleSheet, Alert, ScrollView } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import StyledInput from "../components/StyledInput";
import StyledButton from "../components/StyledButton";
import { COLORS } from "../theme";
import apiClient from "../api/client";
import { AuthContext } from "../context/AuthContext"; // 2. Import AuthContext

const EditProfileScreen = ({ route, navigation }) => {
	const { userRole } = useContext(AuthContext); // 3. Get the user's role
	const existingProfile = route.params?.profile || {};

	// Common fields
	const [firstName, setFirstName] = useState(existingProfile.firstName || "");
	const [lastName, setLastName] = useState(existingProfile.lastName || "");
	// Patient fields
	const [condition, setCondition] = useState(existingProfile.condition || "");
	const [goals, setGoals] = useState(existingProfile.goals || "");
	const [age, setAge] = useState(existingProfile.age ? String(existingProfile.age) : "");
	const [gender, setGender] = useState(existingProfile.gender || "");
	// PT fields
	const [specialty, setSpecialty] = useState(existingProfile.specialty || "");
	const [bio, setBio] = useState(existingProfile.bio || "");

	const [isLoading, setIsLoading] = useState(false);

	const handleSaveProfile = async () => {
		if (!firstName || !lastName) {
			return Alert.alert("Required", "First and last name are required.");
		}
		setIsLoading(true);
		try {
			// 4. Construct the payload based on the user's role
			let profileData = { firstName, lastName };
			if (userRole === "pt") {
				profileData = { ...profileData, specialty, bio };
			} else {
				profileData = { ...profileData, condition, goals, age: age ? Number(age) : undefined, gender };
			}

			await apiClient.post("/profile", profileData);
			Alert.alert("Success", "Your profile has been saved.");
			navigation.goBack();
		} catch (error) {
			Alert.alert("Error", "Could not save profile.");
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<SafeAreaView style={styles.container}>
			<ScrollView>
				<Text style={styles.title}>
					{existingProfile.firstName ? "Edit Profile" : "Create Profile"}
				</Text>
				<StyledInput
					placeholder="First Name"
					value={firstName}
					onChangeText={setFirstName}
				/>
				<StyledInput
					placeholder="Last Name"
					value={lastName}
					onChangeText={setLastName}
				/>

				{/* 5. Conditionally render the correct fields */}
				{userRole === "pt" ? (
					<>
						<StyledInput
							placeholder="Your Specialty (e.g., Sports Injuries)"
							value={specialty}
							onChangeText={setSpecialty}
						/>
						<StyledInput
							placeholder="A short bio about you"
							value={bio}
							onChangeText={setBio}
							multiline
						/>
					</>
				) : (
					<>
						<StyledInput
							placeholder="Primary Condition (e.g., Lower Back Pain)"
							value={condition}
							onChangeText={setCondition}
						/>
						<StyledInput
							placeholder="Your Goals (e.g., Increase mobility)"
							value={goals}
							onChangeText={setGoals}
						/>
					</>
				)}

				<StyledButton
					title="Save Profile"
					onPress={handleSaveProfile}
					isLoading={isLoading}
				/>
			</ScrollView>
		</SafeAreaView>
	);
};

const styles = StyleSheet.create({
	container: { flex: 1, backgroundColor: COLORS.background, padding: 20 },
	title: {
		fontSize: 32,
		fontFamily: "Poppins_700Bold",
		color: COLORS.textDark,
		textAlign: "center",
		marginBottom: 20,
	},
});

export default EditProfileScreen;
