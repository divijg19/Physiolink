// src/models/Appointment.js
const mongoose = require("mongoose");

const AppointmentSchema = new mongoose.Schema({
	pt: {
		// The Physiotherapist offering the slot
		type: mongoose.Schema.Types.ObjectId,
		ref: "user",
		required: true,
	},
	patient: {
		// The Patient who booked the slot
		type: mongoose.Schema.Types.ObjectId,
		ref: "user",
		default: null, // Null until a patient books it
	},
	startTime: {
		type: Date,
		required: true,
	},
	endTime: {
		type: Date,
		required: true,
	},
	status: {
		type: String,
		required: true,
		enum: ["available", "booked", "confirmed", "rejected", "cancelled"],
		default: "available",
	},
	date: {
		type: Date,
		default: Date.now,
	},
});

// To prevent a PT from creating duplicate slots
AppointmentSchema.index({ pt: 1, startTime: 1 }, { unique: true });

module.exports = mongoose.model("appointment", AppointmentSchema);
