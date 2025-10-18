// src/models/User.js
const mongoose = require("mongoose");

// This is the blueprint for a user in our database
const UserSchema = new mongoose.Schema({
	email: {
		type: String,
		required: true,
		unique: true, // No two users can have the same email
	},
	password: {
		type: String,
		required: true, // A password is required
	},
	role: {
		type: String,
		required: true,
		enum: ["patient", "pt", "admin"], // Role must be one of these values
	},
	profile: {
		type: mongoose.Schema.Types.ObjectId,
		ref: "profile",
	},
	date: {
		type: Date,
		default: Date.now, // Automatically sets the registration date
	},
});

module.exports = mongoose.model("user", UserSchema);
