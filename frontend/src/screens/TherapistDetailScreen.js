// src/screens/TherapistDetailScreen.js
import React, { useState, useEffect, useContext } from "react";
import {
    View,
    Text,
    StyleSheet,
    ActivityIndicator,
    FlatList,
} from "react-native";
import Avatar from '../components/Avatar';
import { SafeAreaView } from "react-native-safe-area-context";
import apiClient from "../api/client";
import { Calendar } from 'react-native-calendars';
import { AuthContext } from "../context/AuthContext";
import { COLORS, SPACING, FONT } from "../theme";
import StyledButton from "../components/StyledButton";
import ScreenHeader from "../components/ScreenHeader";
import Snackbar from '../components/Snackbar';
import { formatTime as utilFormatTime, formatDateLong } from '../utils/date';

const formatTime = utilFormatTime;

const TherapistDetailScreen = ({ route, navigation }) => {
    const { therapistId } = route.params;
    const [therapist, setTherapist] = useState(null);
    const [availability, setAvailability] = useState([]);
    const [selectedDay, setSelectedDay] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [bookingSlotId, setBookingSlotId] = useState(null); // To show loading on a specific button
    const [refreshing, setRefreshing] = useState(false);
    const [snack, setSnack] = useState({ visible: false, message: '' });
    const [reviews, setReviews] = useState([]);
    const [myRating, setMyRating] = useState('5');
    const [myComment, setMyComment] = useState('');
    const [canReview, setCanReview] = useState(false);

    useEffect(() => {
        const fetchDetails = async () => {
            try {
                const res = await apiClient.get(`/therapists/${therapistId}`);
                setTherapist(res.data);
                setAvailability(res.data.availableSlots || []);
                // fetch reviews
                const rev = await apiClient.get(`/reviews/${therapistId}`);
                setReviews(rev.data || []);

                // determine if current user (patient) can review: must have a confirmed appointment
                if (userRole !== 'pt') {
                    try {
                        const apps = await apiClient.get('/appointments/me');
                        const myApps = apps.data || [];
                        const hadConfirmed = myApps.some(a => (a.pt?._id || a.pt)?.toString() === therapistId && a.status === 'confirmed');
                        setCanReview(!!hadConfirmed);
                    } catch (err) {
                        // ignore — we'll hide the form
                        setCanReview(false);
                    }
                }
            } catch (error) {
                console.error("Failed to fetch details:", error);
                setSnack({ visible: true, message: 'Could not load therapist details.' });
            } finally {
                setIsLoading(false);
            }
        };
        fetchDetails();
    }, [therapistId]);

    const { userRole } = useContext(AuthContext);

    const handleBookSlot = async (slotId) => {
        setBookingSlotId(slotId); // Show loading spinner on the pressed button
        try {
            await apiClient.put(`/appointments/${slotId}/book`);
            // Remove the booked slot from the local availability list so UI updates immediately
            setAvailability(prev => prev.filter(s => s._id !== slotId));
            setSnack({ visible: true, message: 'Booked — therapist will confirm shortly.' });
            setTimeout(() => {
                const parent = userRole === "pt" ? "TherapistMain" : "PatientMain";
                navigation.navigate(parent, { screen: "Schedule" });
            }, 900);
        } catch (error) {
            const status = error.response?.status;
            const errorMsg = error.response?.data?.msg || "This slot could not be booked.";
            if (status === 409) {
                setSnack({ visible: true, message: "Slot taken — refreshing availability." });
                setTimeout(() => refreshAvailability(), 800);
            } else if (status === 404) {
                setSnack({ visible: true, message: "Slot not found — refreshing availability." });
                setTimeout(() => refreshAvailability(), 800);
            } else {
                setSnack({ visible: true, message: errorMsg || 'Booking failed.' });
            }
        } finally {
            setBookingSlotId(null);
        }
    };

    const refreshAvailability = async () => {
        setRefreshing(true);
        try {
            const res = await apiClient.get(`/therapists/${therapistId}`);
            setAvailability(res.data.availableSlots || []);
        } catch (err) {
            console.error('Failed to refresh availability', err);
        } finally {
            setRefreshing(false);
        }
    };

    if (isLoading || !therapist) {
        return (
            <View style={styles.centered}>
                <ActivityIndicator size="large" color={COLORS.primary} />
            </View>
        );
    }

    const renderSlot = ({ item }) => (
        <View style={styles.slotCard}>
            <Text style={styles.slotText}>
                {formatTime(item.startTime)} - {formatTime(item.endTime)}
            </Text>
            <StyledButton
                title="Book"
                onPress={() => handleBookSlot(item._id)}
                isLoading={bookingSlotId === item._id}
            />
        </View>
    );

    return (
        <SafeAreaView style={styles.container}>
            {/* The header is now provided by the navigator, so we focus on the content */}
            <FlatList
                ListHeaderComponent={
                    <>
                        <Text style={styles.name}>
                            {therapist.profile.firstName} {therapist.profile.lastName}
                        </Text>
                        <Avatar uri={therapist.profile.profileImageUrl} name={`${therapist.profile.firstName} ${therapist.profile.lastName}`} size={96} />
                        <Text style={styles.specialty}>{therapist.profile.specialty}</Text>
                        {/* Show rating from profile if present */}
                        <Text style={styles.rating}>
                            ⭐ {therapist.profile?.rating?.toFixed(1) || 'N/A'}{' '}
                            {therapist.reviewCount ? `• ${therapist.reviewCount} ${therapist.reviewCount === 1 ? 'review' : 'reviews'}` : '• No reviews'}
                        </Text>
                        {/* Calendar to select a day */}
                        <Calendar
                            onDayPress={(day) => setSelectedDay(day.dateString)}
                            markedDates={selectedDay ? { [selectedDay]: { selected: true } } : {}}
                            style={{ marginVertical: SPACING.sm }}
                        />
                        <View style={[styles.infoBox, { padding: SPACING.md }]}>
                            <Text style={styles.bio}>
                                {therapist.profile.bio || "No biography provided."}
                            </Text>
                            {therapist.profile.credentials ? (
                                <Text style={{ marginTop: 8, fontFamily: 'Poppins_400Regular' }}>
                                    Credentials: {therapist.profile.credentials}
                                </Text>
                            ) : null}
                            {therapist.profile.location ? (
                                <Text style={{ marginTop: 8, fontFamily: 'Poppins_400Regular' }}>
                                    Location: {therapist.profile.location}
                                </Text>
                            ) : null}
                            <View style={{ marginTop: 12 }}>
                                <Text style={{ fontFamily: 'Poppins_600SemiBold' }}>Reviews</Text>
                                {reviews.length === 0 ? (
                                    <Text style={{ color: COLORS.gray }}>No reviews yet.</Text>
                                ) : (
                                    reviews.slice(0, 3).map(r => (
                                        <View key={r._id} style={{ marginTop: 8 }}>
                                            <Text style={{ fontFamily: 'Poppins_600SemiBold' }}>{r.patient?.profile?.firstName || 'Patient'} • ⭐ {r.rating}</Text>
                                            {r.comment ? <Text style={{ fontFamily: 'Poppins_400Regular' }}>{r.comment}</Text> : null}
                                        </View>
                                    ))
                                )}
                            </View>
                            {userRole !== 'pt' && canReview && (
                                <View style={{ marginTop: 12 }}>
                                    <Text style={{ fontFamily: 'Poppins_600SemiBold' }}>Leave a review</Text>
                                    <View style={{ flexDirection: 'row', marginTop: 8, alignItems: 'center' }}>
                                        <Text style={{ marginRight: 8 }}>Rating:</Text>
                                        <StyledInput value={myRating} onChangeText={setMyRating} style={{ width: 60 }} />
                                    </View>
                                    <StyledInput placeholder="Comment (optional)" value={myComment} onChangeText={setMyComment} multiline />
                                    <StyledButton title="Submit Review" onPress={async () => {
                                        try {
                                            const payload = { therapistId, rating: Number(myRating), comment: myComment };
                                            await apiClient.post('/reviews', payload);
                                            setSnack({ visible: true, message: 'Thank you — review submitted.' });
                                            // refresh reviews and therapist rating
                                            const rev = await apiClient.get(`/reviews/${therapistId}`);
                                            setReviews(rev.data || []);
                                            const t = await apiClient.get(`/therapists/${therapistId}`);
                                            setTherapist(t.data);
                                        } catch (err) {
                                            console.error('Failed to submit review', err);
                                            setSnack({ visible: true, message: 'Could not submit review.' });
                                        }
                                    }} />
                                </View>
                            )}
                        </View>
                        <ScreenHeader title="Available Slots" />
                    </>
                }
                data={
                    selectedDay
                        ? availability.filter(a => new Date(a.startTime).toISOString().slice(0, 10) === selectedDay)
                        : availability
                }
                renderItem={renderSlot}
                keyExtractor={(item) => item._id}
                ListEmptyComponent={
                    <Text style={styles.placeholderText}>
                        This therapist has no available slots.
                    </Text>
                }
                contentContainerStyle={{ paddingHorizontal: 20 }}
            />
            <Snackbar message={snack.message} visible={snack.visible} />
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    centered: { flex: 1, justifyContent: "center", alignItems: "center" },
    container: { flex: 1, backgroundColor: COLORS.background },
    name: {
        fontFamily: "Poppins_700Bold",
        fontSize: 28,
        color: COLORS.textDark,
        textAlign: "center",
        marginTop: 10,
    },
    specialty: {
        fontFamily: "Poppins_400Regular",
        fontSize: 16,
        color: COLORS.gray,
        textAlign: "center",
        marginBottom: 20,
    },
    infoBox: {
        backgroundColor: COLORS.lightGray,
        padding: 20,
        borderRadius: 12,
        marginBottom: 20,
    },
    bio: {
        fontFamily: "Poppins_400Regular",
        fontSize: 15,
        color: COLORS.textDark,
        lineHeight: 22,
    },
    placeholderText: {
        fontFamily: "Poppins_400Regular",
        color: COLORS.gray,
        textAlign: "center",
        marginVertical: 20,
    },
    slotCard: {
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "center",
        backgroundColor: COLORS.lightGray,
        padding: 10,
        borderRadius: 10,
        marginVertical: 5,
    },
    slotText: {
        fontFamily: "Poppins_600SemiBold",
        fontSize: 16,
        color: COLORS.textDark,
        marginLeft: 10,
    },
});

export default TherapistDetailScreen;
