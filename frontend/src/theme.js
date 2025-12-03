// src/theme.js
export const COLORS = {
	primary: "#0057FF", // A strong, trustworthy blue
	secondary: "#00A3FF", // A lighter, friendly blue
	accent: "#FF7A59", // A soft, encouraging orange for buttons

	textDark: "#1E293B", // A dark gray, not pure black, easier on the eyes
	textLight: "#F8FAFC",

	background: "#FFFFFF",
	gray: "#E2E8F0",
	lightGray: "#F1F5F9",
	darkGray: '#2D3748',
};

// Spacing, typography and subtle elevation tokens for consistent spacing and a premium feel
export const SPACING = {
	xs: 6,
	sm: 10,
	md: 16,
	lg: 24,
	xl: 32,
};

export const FONT = {
	largeTitle: 28,
	title: 20,
	body: 16,
	small: 13,
};

export const SHADOW = {
	// cross-platform friendly shadow preset
	elevation: 3,
	shadowColor: '#000',
	shadowOffset: { width: 0, height: 2 },
	shadowOpacity: 0.08,
	shadowRadius: 6,
};
