// src/screens/TherapistDetailScreen.js
import React, { useState, useEffect, useContext } from "react";
import {
    View,
    Text,
    StyleSheet,
    ActivityIndicator,
    Alert,
    FlatList,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import apiClient from "../api/client";
import { Calendar } from 'react-native-calendars';
import { AuthContext } from "../context/AuthContext";
import { COLORS } from "../theme";
import StyledButton from "../components/StyledButton";
import ScreenHeader from "../components/ScreenHeader";

// Re-using our date/time formatters from MyScheduleScreen
const formatTime = (timeString) => {
    const options = { hour: "numeric", minute: "2-digit" };
    return new Date(timeString).toLocaleTimeString(undefined, options);
};

const TherapistDetailScreen = ({ route, navigation }) => {
    const { therapistId } = route.params;
    const [therapist, setTherapist] = useState(null);
    const [availability, setAvailability] = useState([]);
    const [selectedDay, setSelectedDay] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [bookingSlotId, setBookingSlotId] = useState(null); // To show loading on a specific button

    useEffect(() => {
        const fetchDetails = async () => {
            try {
                // Fetch both therapist profile and availability in parallel
                const [therapistsRes, availabilityRes] = await Promise.all([
                    apiClient.get("/therapists"),
                    apiClient.get(`/appointments/availability/${therapistId}`),
                ]);

                const foundTherapist = therapistsRes.data.find((t) => t._id === therapistId);
                setTherapist(foundTherapist);
                setAvailability(availabilityRes.data);
            } catch (error) {
                console.error("Failed to fetch details:", error);
                Alert.alert("Error", "Could not load therapist details.");
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
            Alert.alert(
                "Success!",
                "Your appointment has been booked. The therapist will confirm shortly.",
                [
                    {
                        text: "OK",
                        onPress: () => {
                            // Navigate to the tab navigator's Schedule screen. The AppStack nests either
                            // 'PatientMain' or 'TherapistMain' depending on role. We route accordingly.
                            const parent =
                                userRole === "pt" ? "TherapistMain" : "PatientMain";
                            navigation.navigate(parent, { screen: "Schedule" });
                        },
                    },
                ],
            );
        } catch (error) {
            const errorMsg =
                error.response?.data?.msg || "This slot could not be booked.";
            Alert.alert("Booking Failed", errorMsg);
        } finally {
            setBookingSlotId(null);
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
                        <Text style={styles.specialty}>{therapist.profile.specialty}</Text>
                        {/* Static rating for MVP */}
                        <Text style={styles.rating}>⭐ 4.5</Text>
                        {/* Calendar to select a day */}
                        <Calendar
                            onDayPress={(day) => setSelectedDay(day.dateString)}
                            markedDates={selectedDay ? { [selectedDay]: { selected: true } } : {}}
                            style={{ marginVertical: 10 }}
                        />
                        <View style={styles.infoBox}>
                            <Text style={styles.bio}>
                                {therapist.profile.bio || "No biography provided."}
                            </Text>
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
