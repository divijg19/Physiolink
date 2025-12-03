// src/models/Profile.js
const mongoose = require("mongoose");

const ProfileSchema = new mongoose.Schema({
	// Link to the User model
	user: {
		type: mongoose.Schema.Types.ObjectId,
		ref: "user",
		required: true,
	},
	// --- Common Fields ---
	firstName: {
		type: String,
		required: true,
	},
	lastName: {
		type: String,
		required: true,
	},
	// --- Patient-Specific Fields ---
	age: {
		type: Number,
	},
	gender: {
		type: String,
	},
	condition: {
		// e.g., "Lower Back Pain"
		type: String,
	},
	goals: {
		// e.g., "Increase mobility"
		type: String,
	},
	// --- PT-Specific Fields ---
	specialty: {
		// e.g., "Sports Injury", "Pediatrics"
		type: String,
	},
	location: {
		// e.g., "San Francisco, CA"
		type: String,
	},
	bio: {
		type: String,
	},
	credentials: {
		// URL to a document, or text
		type: String,
	},
	profileImageUrl: {
		// URL to a profile image
		type: String,
	},
	rating: {
		type: Number,
		min: 0,
		max: 5,
		default: 0,
	},
	isVerified: {
		type: Boolean,
		default: false,
	},
	date: {
		type: Date,
		default: Date.now,
	},
});

module.exports = mongoose.model("profile", ProfileSchema);
