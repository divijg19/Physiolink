// src/config/db.js
const mongoose = require("mongoose");
require("dotenv").config({ path: "../../.env" }); // Adjust path to find .env

const connectDB = async () => {
	try {
		await mongoose.connect(process.env.MONGO_URI);
		console.log("MongoDB Connected... Success!");
	} catch (err) {
		console.error(err);
		// Exit process with failure
		process.exit(1);
	}
};

module.exports = connectDB;
