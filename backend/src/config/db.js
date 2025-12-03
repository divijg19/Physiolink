// src/config/db.js
const mongoose = require("mongoose");
require("dotenv").config({ path: "../../.env" }); // Adjust path to find .env

const connectDB = async () => {
	try {
		await mongoose.connect(process.env.MONGO_URI);
		console.log("MongoDB Connected... Success!");
	} catch (err) {
		console.error(err);
		// In production we exit, but during tests let the error bubble so Jest can handle it
		if (process.env.NODE_ENV === "test") {
			throw err;
		}
		// Exit process with failure for non-test environments
		process.exit(1);
	}
};

module.exports = connectDB;
