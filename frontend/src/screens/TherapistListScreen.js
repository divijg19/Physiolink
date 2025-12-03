// src/screens/TherapistListScreen.js
import React, { useState, useEffect, useCallback } from "react";
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
import { COLORS, SPACING, FONT, SHADOW } from "../theme";
import ScreenHeader from "../components/ScreenHeader";
import StyledInput from "../components/StyledInput";
import Avatar from '../components/Avatar';

const TherapistListScreen = ({ navigation }) => {
	const [therapists, setTherapists] = useState([]);
	const [isLoading, setIsLoading] = useState(true);
	const [isLoadingMore, setIsLoadingMore] = useState(false);
	const [page, setPage] = useState(1);
	const [totalPages, setTotalPages] = useState(1);
	const [specialtyFilter, setSpecialtyFilter] = useState("");
	const [locationFilter, setLocationFilter] = useState("");
	const [debounceTimer, setDebounceTimer] = useState(null);
	const [availableOnly, setAvailableOnly] = useState(false);

	const fetchTherapists = useCallback(async (opts = {}) => {
		const { next = false, reset = false } = opts;
		try {
			if (reset) {
				setIsLoading(true);
				setPage(1);
			}
			if (next) setIsLoadingMore(true);
			const params = { page: next ? page + 1 : page, limit: 10 };
			if (specialtyFilter) params.specialty = specialtyFilter;
			if (locationFilter) params.location = locationFilter;
			if (availableOnly) params.available = true;
			const response = await apiClient.get('/therapists', { params });
			// response shape: { data, total, page, totalPages }
			const { data, page: respPage, totalPages: respTotalPages } = response.data;
			const filtered = (data || []).filter((t) => t.profile);
			if (reset) {
				setTherapists(filtered);
			} else if (next) {
				setTherapists(prev => [...prev, ...filtered]);
			} else {
				setTherapists(filtered);
			}
			setPage(respPage || params.page);
			setTotalPages(respTotalPages || 1);
		} catch (error) {
			console.error('Failed to fetch therapists:', error);
		} finally {
			setIsLoading(false);
			setIsLoadingMore(false);
		}
	}, [page, specialtyFilter, locationFilter, availableOnly]);

	// when filters change, debounce to avoid excessive requests
	useEffect(() => {
		if (debounceTimer) clearTimeout(debounceTimer);
		const t = setTimeout(() => fetchTherapists({ reset: true }), 400);
		setDebounceTimer(t);
		return () => clearTimeout(t);
	}, [specialtyFilter, locationFilter, fetchTherapists]);

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
			<View style={styles.row}>
				<Avatar uri={item.profile.profileImageUrl} name={`${item.profile.firstName} ${item.profile.lastName}`} size={56} />
				<View style={styles.meta}>
					<Text style={styles.cardTitle}>
						{item.profile.firstName} {item.profile.lastName}
					</Text>
					<Text style={styles.cardSubtitle}>
						{item.profile.specialty || "General Physiotherapy"}
					</Text>
				</View>
				<View style={styles.badgeCol}>
					<Text style={styles.smallMuted}>
						{item.profile.rating ? `⭐ ${item.profile.rating.toFixed(1)} • ${item.reviewCount || 0} ${item.reviewCount === 1 ? 'review' : 'reviews'}` : `No reviews`}
					</Text>
					<Text style={styles.slotsText}>{item.availableSlotsCount ? `${item.availableSlotsCount} slots` : 'No slots'}</Text>
				</View>
			</View>
		</TouchableOpacity>
	);

	return (
		<SafeAreaView style={styles.container}>
			<ScreenHeader
				title="Find a Therapist"
				subtitle="Browse available specialists"
			/>
			<View style={styles.filtersWrap}>
				<StyledInput placeholder="Search by Specialty" value={specialtyFilter} onChangeText={setSpecialtyFilter} />
				<StyledInput placeholder="Search by Location" value={locationFilter} onChangeText={setLocationFilter} />
				<TouchableOpacity style={[styles.availabilityToggle, availableOnly && styles.availabilityToggleActive]} onPress={() => { setAvailableOnly(v => !v); fetchTherapists({ reset: true }); }}>
					<Text style={availableOnly ? styles.toggleActiveText : styles.toggleText}>{availableOnly ? 'Showing: Available only' : 'Show only available'}</Text>
				</TouchableOpacity>
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
			{page < totalPages && (
				<TouchableOpacity
					style={styles.loadMore}
					onPress={() => fetchTherapists({ next: true })}
					disabled={isLoadingMore}
				>
					{isLoadingMore ? (
						<ActivityIndicator color={COLORS.primary} />
					) : (
						<Text style={styles.loadMoreText}>Load more</Text>
					)}
				</TouchableOpacity>
			)}
		</SafeAreaView>
	);
};

const styles = StyleSheet.create({
	centered: { flex: 1, justifyContent: "center", alignItems: "center" },
	container: { flex: 1, backgroundColor: COLORS.background },
	card: {
		backgroundColor: COLORS.lightGray,
		padding: SPACING.md,
		marginVertical: SPACING.sm,
		marginHorizontal: SPACING.md,
		borderRadius: 12,
		...SHADOW,
	},
	row: { flexDirection: 'row', alignItems: 'center' },
	meta: { marginLeft: SPACING.sm, flex: 1 },
	badgeCol: { alignItems: 'flex-end' },
	cardTitle: {
		fontFamily: "Poppins_600SemiBold",
		fontSize: FONT.title,
		color: COLORS.textDark,
	},
	cardSubtitle: {
		fontFamily: "Poppins_400Regular",
		fontSize: FONT.body,
		color: COLORS.gray,
		marginTop: SPACING.xs,
	},
	smallMuted: {
		color: COLORS.gray,
		fontSize: FONT.small,
		fontFamily: 'Poppins_400Regular',
	},
	slotsText: { color: COLORS.primary, fontSize: FONT.small, fontFamily: 'Poppins_600SemiBold', marginTop: SPACING.xs },
	emptyText: {
		textAlign: "center",
		marginTop: SPACING.xl,
		fontFamily: "Poppins_400Regular",
		color: COLORS.gray,
	},
	loadMore: {
		alignSelf: 'center',
		marginVertical: SPACING.md,
		backgroundColor: COLORS.primary,
		paddingHorizontal: SPACING.lg,
		paddingVertical: SPACING.sm,
		borderRadius: 8,
	},
	loadMoreText: {
		color: '#fff',
		fontFamily: 'Poppins_600SemiBold',
	},
	filtersWrap: { paddingHorizontal: SPACING.md, paddingBottom: SPACING.sm },
	availabilityToggle: {
		marginTop: SPACING.sm,
		padding: SPACING.sm,
		borderRadius: 8,
		alignItems: 'center',
		borderWidth: 1,
		borderColor: COLORS.primary,
	},
	availabilityToggleActive: {
		backgroundColor: COLORS.primary,
	},
	toggleText: { color: COLORS.primary },
	toggleActiveText: { color: COLORS.white },
});

export default TherapistListScreen;
