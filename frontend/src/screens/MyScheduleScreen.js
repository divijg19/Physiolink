// src/screens/MyScheduleScreen.js
import React, { useCallback, useState, useContext } from "react";
import {
    View,
    Text,
    StyleSheet,
    FlatList,
    ActivityIndicator,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useFocusEffect } from "@react-navigation/native";
import { COLORS } from "../theme";
import apiClient from "../api/client";
import ScreenHeader from "../components/ScreenHeader";
import { AuthContext } from "../context/AuthContext";

// A simple utility to format dates nicely
const formatDate = (dateString) => {
    const options = { weekday: "long", month: "long", day: "numeric" };
    return new Date(dateString).toLocaleDateString(undefined, options);
};

// A simple utility to format times nicely
const formatTime = (timeString) => {
    const options = { hour: "numeric", minute: "2-digit", hour12: true };
    return new Date(timeString).toLocaleTimeString(undefined, options);
};

// The AppointmentCard is now smart enough to display the correct person's name
const AppointmentCard = ({ item, userRole, onStatusChange }) => {
    // Determine who the other person in the appointment is
    const otherParty = userRole === "pt" ? item.patient : item.pt;

    // A robust check to ensure the other user's profile exists before trying to display a name
    const otherPartyName = otherParty?.profile
        ? `${otherParty.profile.firstName} ${otherParty.profile.lastName}`
        : "Unassigned"; // Fallback text

    return (
        <View style={styles.card}>
            <Text style={styles.cardWith}>
                {userRole === "pt" ? "With Patient:" : "With Therapist:"}
            </Text>
            <Text style={styles.cardName}>{otherPartyName}</Text>
            <View style={styles.timeInfo}>
                <Text style={styles.cardDate}>{formatDate(item.startTime)}</Text>
                <Text style={styles.cardTime}>
                    {formatTime(item.startTime)} - {formatTime(item.endTime)}
                </Text>
            </View>
            <Text style={styles.cardStatus(item.status)}>
                {item.status.toUpperCase()}
            </Text>
            {userRole === 'pt' && item.status === 'booked' && (
                <View style={{ marginTop: 12 }}>
                    <StyledButton title="Accept" onPress={() => onStatusChange(item._id, 'confirmed')} />
                    <View style={{ height: 8 }} />
                    <StyledButton title="Reject" onPress={() => onStatusChange(item._id, 'rejected')} />
                </View>
            )}
        </View>
    );
};

const MyScheduleScreen = () => {
    const { userRole } = useContext(AuthContext); // Get the user's role
    const [appointments, setAppointments] = useState([]);
    const [isLoading, setIsLoading] = useState(true);

    const handleStatusChange = async (id, status) => {
        try {
            await apiClient.put(`/appointments/${id}/status`, { status });
            // Refresh schedule
            const response = await apiClient.get('/appointments/me');
            setAppointments(response.data);
        } catch (error) {
            console.error('Could not update status', error);
        }
    };

    // useFocusEffect re-fetches data every time the tab is viewed, so it's always up-to-date
    useFocusEffect(
        useCallback(() => {
            const fetchSchedule = async () => {
                if (!isLoading) setIsLoading(true); // Show loading indicator on re-focus
                try {
                    const response = await apiClient.get("/appointments/me");
                    setAppointments(response.data);
                } catch (error) {
                    console.error("Failed to fetch schedule:", error);
                } finally {
                    setIsLoading(false);
                }
            };
            fetchSchedule();
        }, []),
    );

    return (
        <SafeAreaView style={styles.container}>
            <ScreenHeader title="My Schedule" subtitle="Your upcoming appointments" />
            {isLoading ? (
                <ActivityIndicator
                    size="large"
                    color={COLORS.primary}
                    style={{ marginTop: 50 }}
                />
            ) : (
                <FlatList
                    data={appointments}
                    renderItem={({ item }) => (
                        <AppointmentCard item={item} userRole={userRole} onStatusChange={handleStatusChange} />
                    )}
                    keyExtractor={(item) => item._id}
                    contentContainerStyle={{ paddingHorizontal: 16, paddingBottom: 20 }}
                    ListEmptyComponent={
                        <View style={styles.emptyContainer}>
                            <Text style={styles.emptyText}>
                                You have no appointments scheduled.
                            </Text>
                        </View>
                    }
                />
            )}
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: { flex: 1, backgroundColor: COLORS.background },
    emptyContainer: {
        alignItems: "center",
        marginTop: 50,
        paddingHorizontal: 20,
    },
    emptyText: {
        fontFamily: "Poppins_400Regular",
        color: COLORS.gray,
        fontSize: 16,
        textAlign: "center",
    },
    card: {
        backgroundColor: COLORS.lightGray,
        padding: 20,
        marginVertical: 8,
        borderRadius: 12,
    },
    cardWith: {
        fontFamily: "Poppins_400Regular",
        fontSize: 13,
        color: COLORS.gray,
    },
    cardName: {
        fontFamily: "Poppins_600SemiBold",
        fontSize: 18,
        color: COLORS.textDark,
        marginBottom: 10,
    },
    timeInfo: {
        borderTopWidth: 1,
        borderTopColor: "#E2E8F0",
        paddingTop: 10,
        marginTop: 5,
    },
    cardDate: {
        fontFamily: "Poppins_600SemiBold",
        fontSize: 15,
        color: COLORS.textDark,
    },
    cardTime: {
        fontFamily: "Poppins_400Regular",
        fontSize: 14,
        color: COLORS.textDark,
        marginTop: 4,
    },
    cardStatus: (status) => ({
        fontFamily: "Poppins_600SemiBold",
        fontSize: 12,
        color:
            status === "available"
                ? COLORS.primary
                : status === "booked"
                    ? COLORS.accent
                    : COLORS.gray,
        position: "absolute",
        top: 20,
        right: 20,
        textTransform: "uppercase",
    }),
});

export default MyScheduleScreen;
