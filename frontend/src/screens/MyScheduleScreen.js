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
import { COLORS, SPACING, FONT, SHADOW } from "../theme";
import apiClient from "../api/client";
import ScreenHeader from "../components/ScreenHeader";
import StyledButton from "../components/StyledButton";
import Snackbar from '../components/Snackbar';
import { AuthContext } from "../context/AuthContext";
import { formatDateLong, formatTime } from '../utils/date';

// A simple utility to format dates nicely
// Use shared utilities
const formatDate = formatDateLong;

// Appointment card with optional loading indicators for PT actions
const AppointmentCard = ({ item, userRole, onStatusChange, loadingIds = {} }) => {
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
                    <StyledButton title="Accept" onPress={() => onStatusChange(item._id, 'confirmed')} isLoading={!!loadingIds[item._id]} />
                    <View style={{ height: 8 }} />
                    <StyledButton title="Reject" onPress={() => onStatusChange(item._id, 'rejected')} isLoading={!!loadingIds[item._id]} />
                </View>
            )}
        </View>
    );
};

const MyScheduleScreen = () => {
    const { userRole } = useContext(AuthContext); // Get the user's role
    const [appointments, setAppointments] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [reminders, setReminders] = useState([]);
    const [actionLoading, setActionLoading] = useState({});
    const [snack, setSnack] = useState({ visible: false, message: '' });

    const handleStatusChange = async (id, status) => {
        setActionLoading(prev => ({ ...prev, [id]: true }));
        try {
            await apiClient.put(`/appointments/${id}/status`, { status });
            // Refresh schedule
            const response = await apiClient.get('/appointments/me');
            setAppointments(response.data);
            setSnack({ visible: true, message: `Appointment ${status}` });
        } catch (error) {
            console.error('Could not update status', error);
            setSnack({ visible: true, message: 'Action failed — please try again.' });
        } finally {
            setActionLoading(prev => ({ ...prev, [id]: false }));
        }
    };

    // useFocusEffect re-fetches data every time the tab is viewed, so it's always up-to-date
    useFocusEffect(
        useCallback(() => {
            const fetchSchedule = async () => {
                if (!isLoading) setIsLoading(true); // Show loading indicator on re-focus
                try {
                    const [appsRes, remindersRes] = await Promise.all([
                        apiClient.get('/appointments/me'),
                        apiClient.get('/reminders/me').catch(() => ({ data: [] })),
                    ]);
                    setAppointments(appsRes.data);
                    setReminders(remindersRes.data || []);
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
                    style={styles.loadingIndicator}
                />
            ) : (
                <>
                    {reminders.length > 0 && (
                        <View style={styles.remindersWrap}>
                            <Text style={styles.remindersTitle}>Reminders</Text>
                            {reminders.map(r => (
                                <View key={r._id} style={styles.reminderItem}>
                                    <Text style={styles.reminderText}>{r.message}</Text>
                                    <Text style={styles.reminderTime}>{new Date(r.remindAt).toLocaleString()}</Text>
                                </View>
                            ))}
                        </View>
                    )}
                    <FlatList
                        data={appointments}
                        renderItem={({ item }) => (
                            <AppointmentCard item={item} userRole={userRole} onStatusChange={handleStatusChange} loadingIds={actionLoading} />
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
                    <Snackbar message={snack.message} visible={snack.visible} />
                </>
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
        padding: SPACING.md,
        marginVertical: SPACING.sm,
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
        marginBottom: SPACING.sm,
    },
    timeInfo: {
        borderTopWidth: 1,
        borderTopColor: "#E2E8F0",
        paddingTop: SPACING.sm,
        marginTop: SPACING.xs,
    },
    cardDate: {
        fontFamily: "Poppins_600SemiBold",
        fontSize: FONT.body,
        color: COLORS.textDark,
    },
    cardTime: {
        fontFamily: "Poppins_400Regular",
        fontSize: FONT.body,
        color: COLORS.textDark,
        marginTop: SPACING.xs,
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
    loadingIndicator: { marginTop: SPACING.lg },
    remindersWrap: { paddingHorizontal: SPACING.md, paddingVertical: SPACING.sm },
    remindersTitle: { fontFamily: 'Poppins_600SemiBold' },
    reminderItem: { backgroundColor: COLORS.lightGray, padding: SPACING.sm, borderRadius: 8, marginTop: SPACING.sm },
    reminderText: { fontFamily: 'Poppins_400Regular' },
    reminderTime: { fontFamily: 'Poppins_400Regular', color: COLORS.gray, marginTop: SPACING.xs },
    actionRow: { marginTop: SPACING.md },
});

export default MyScheduleScreen;
