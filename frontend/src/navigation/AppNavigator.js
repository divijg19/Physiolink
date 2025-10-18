import React, { useContext } from "react";
import { View, ActivityIndicator, TouchableOpacity } from "react-native";
import { NavigationContainer } from "@react-navigation/native";
import { createStackNavigator } from "@react-navigation/stack";
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import { AuthContext } from "../context/AuthContext";
import { Feather } from "@expo/vector-icons";
import { COLORS } from "../theme";

// --- Screen Imports ---
// Auth Screens
import LoginScreen from "../screens/LoginScreen";
import RegisterScreen from "../screens/RegisterScreen";
// Shared App Screens
import ProfileScreen from "../screens/ProfileScreen";
import EditProfileScreen from "../screens/EditProfileScreen";
import MyScheduleScreen from "../screens/MyScheduleScreen";
// Patient-Only Screens
import TherapistListScreen from "../screens/TherapistListScreen";
import TherapistDetailScreen from "../screens/TherapistDetailScreen";
// Therapist-Only Screens
import SetAvailabilityScreen from "../screens/SetAvailabilityScreen";

const AuthStack = createStackNavigator();
const AppStack = createStackNavigator();
const Tab = createBottomTabNavigator();

// This stack remains the same for the login/register flow
const AuthNavigator = () => (
	<AuthStack.Navigator screenOptions={{ headerShown: false }}>
		<AuthStack.Screen name="Login" component={LoginScreen} />
		<AuthStack.Screen name="Register" component={RegisterScreen} />
	</AuthStack.Navigator>
);

// Define the set of tabs for a PATIENT
const PatientTabs = () => (
	<Tab.Navigator
		screenOptions={({ route }) => ({
			headerShown: false,
			tabBarIcon: ({ focused, color, size }) => {
				let iconName;
				if (route.name === "Discover") iconName = "search";
				else if (route.name === "Schedule") iconName = "calendar";
				else if (route.name === "ProfileTab") iconName = "user";
				return <Feather name={iconName} size={size} color={color} />;
			},
			tabBarActiveTintColor: COLORS.primary,
			tabBarInactiveTintColor: COLORS.gray,
			tabBarStyle: { backgroundColor: COLORS.background },
		})}
	>
		<Tab.Screen name="Discover" component={TherapistListScreen} />
		<Tab.Screen name="Schedule" component={MyScheduleScreen} />
		<Tab.Screen
			name="ProfileTab"
			component={ProfileScreen}
			options={{ title: "Profile" }}
		/>
	</Tab.Navigator>
);

// Define the set of tabs for a THERAPIST
const TherapistTabs = () => (
	<Tab.Navigator
		screenOptions={({ route }) => ({
			headerShown: false,
			tabBarIcon: ({ focused, color, size }) => {
				let iconName;
				if (route.name === "Schedule") iconName = "calendar";
				else if (route.name === "Availability") iconName = "clock";
				else if (route.name === "ProfileTab") iconName = "user";
				return <Feather name={iconName} size={size} color={color} />;
			},
			tabBarActiveTintColor: COLORS.primary,
			tabBarInactiveTintColor: COLORS.gray,
			tabBarStyle: { backgroundColor: COLORS.background },
		})}
	>
		<Tab.Screen name="Schedule" component={MyScheduleScreen} />
		<Tab.Screen name="Availability" component={SetAvailabilityScreen} />
		<Tab.Screen
			name="ProfileTab"
			component={ProfileScreen}
			options={{ title: "Profile" }}
		/>
	</Tab.Navigator>
);

export const AppNavigator = () => {
	const { userToken, userRole, isLoading } = useContext(AuthContext);

	if (isLoading) {
		return (
			<View style={{ flex: 1, justifyContent: "center", alignItems: "center" }}>
				<ActivityIndicator size="large" color={COLORS.primary} />
			</View>
		);
	}

	return (
		<NavigationContainer>
			{userToken ? (
				// REFINEMENT: This AppStack now provides consistent, professional headers
				<AppStack.Navigator
					screenOptions={({ navigation }) => ({
						headerStyle: {
							backgroundColor: COLORS.background,
							elevation: 0, // Remove shadow on Android
							shadowOpacity: 0, // Remove shadow on iOS
						},
						headerTitleStyle: {
							fontFamily: "Poppins_600SemiBold",
							color: COLORS.textDark,
						},
						// Add a custom back button to all screens
						headerLeft: () => (
							<TouchableOpacity
								onPress={() => navigation.goBack()}
								style={{ marginLeft: 16 }}
							>
								<Feather
									name="chevron-left"
									size={28}
									color={COLORS.textDark}
								/>
							</TouchableOpacity>
						),
					})}
				>
					{/* The Tab Navigators are the "base" screens and should NOT show a header */}
					{userRole === "pt" ? (
						<AppStack.Screen
							name="TherapistMain"
							component={TherapistTabs}
							options={{ headerShown: false }}
						/>
					) : (
						<AppStack.Screen
							name="PatientMain"
							component={PatientTabs}
							options={{ headerShown: false }}
						/>
					)}
					{/* Screens pushed ON TOP of the tabs WILL show the header we just defined */}
					<AppStack.Screen
						name="TherapistDetail"
						component={TherapistDetailScreen}
						options={{ title: "Therapist Profile" }}
					/>
					<AppStack.Screen
						name="EditProfile"
						component={EditProfileScreen}
						options={{ title: "Edit Profile" }}
					/>
				</AppStack.Navigator>
			) : (
				<AuthNavigator />
			)}
		</NavigationContainer>
	);
};
