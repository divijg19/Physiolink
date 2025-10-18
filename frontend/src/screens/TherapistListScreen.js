// src/screens/TherapistListScreen.js
import React, { useState, useEffect } from "react";
// THE FIX: Removed SafeAreaView from 'react-native'
import {
	View,
	Text,
	StyleSheet,
	FlatList,
	ActivityIndicator,
	TouchableOpacity,
} from "react-native";
// THE FIX: Added the correct import
import { SafeAreaView } from "react-native-safe-area-context";
import apiClient from "../api/client";
import { COLORS } from "../theme";
import ScreenHeader from "../components/ScreenHeader";
import StyledInput from "../components/StyledInput";

const TherapistListScreen = ({ navigation }) => {
	const [therapists, setTherapists] = useState([]);
	const [isLoading, setIsLoading] = useState(true);
	const [specialtyFilter, setSpecialtyFilter] = useState("");
	const [locationFilter, setLocationFilter] = useState("");

	useEffect(() => {
		const fetchTherapists = async () => {
			setIsLoading(true);
			try {
				const params = {};
				if (specialtyFilter) params.specialty = specialtyFilter;
				if (locationFilter) params.location = locationFilter;
				const response = await apiClient.get('/therapists', { params });
				const therapistsWithProfiles = response.data.filter((t) => t.profile);
				setTherapists(therapistsWithProfiles);
			} catch (error) {
				console.error('Failed to fetch therapists:', error);
			} finally {
				setIsLoading(false);
			}
		};
		fetchTherapists();
	}, [specialtyFilter, locationFilter]);

	if (isLoading) {
		return (
			<View style={styles.centered}>
				<ActivityIndicator size="large" color={COLORS.primary} />
			</View>
		);
	}

	const renderTherapist = ({ item }) => (
		<TouchableOpacity
			style={styles.card}
			onPress={() =>
				navigation.navigate("TherapistDetail", { therapistId: item._id })
			}
		>
			<Text style={styles.cardTitle}>
				{item.profile.firstName} {item.profile.lastName}
			</Text>
			<Text style={styles.cardSubtitle}>
				{item.profile.specialty || "General Physiotherapy"}
			</Text>
		</TouchableOpacity>
	);

	return (
		<SafeAreaView style={styles.container}>
			<ScreenHeader
				title="Find a Therapist"
				subtitle="Browse available specialists"
			/>
			<View style={{ paddingHorizontal: 16 }}>
				<StyledInput placeholder="Search by Specialty" value={specialtyFilter} onChangeText={setSpecialtyFilter} />
				<StyledInput placeholder="Search by Location" value={locationFilter} onChangeText={setLocationFilter} />
			</View>
			<FlatList
				data={therapists}
				renderItem={renderTherapist}
				keyExtractor={(item) => item._id}
				ListEmptyComponent={
					<Text style={styles.emptyText}>
						No therapists available at the moment.
					</Text>
				}
			/>
		</SafeAreaView>
	);
};

const styles = StyleSheet.create({
	centered: { flex: 1, justifyContent: "center", alignItems: "center" },
	container: { flex: 1, backgroundColor: COLORS.background },
	card: {
		backgroundColor: COLORS.lightGray,
		padding: 20,
		marginVertical: 8,
		marginHorizontal: 16,
		borderRadius: 12,
	},
	cardTitle: {
		fontFamily: "Poppins_600SemiBold",
		fontSize: 18,
		color: COLORS.textDark,
	},
	cardSubtitle: {
		fontFamily: "Poppins_400Regular",
		fontSize: 14,
		color: COLORS.gray,
		marginTop: 4,
	},
	emptyText: {
		textAlign: "center",
		marginTop: 50,
		fontFamily: "Poppins_400Regular",
		color: COLORS.gray,
	},
});

export default TherapistListScreen;
