// scripts/seed.js
const _mongoose = require("mongoose");
require("dotenv").config();
const connectDB = require("../src/config/db");
const User = require("../src/models/User");
const Profile = require("../src/models/Profile");
const bcrypt = require("bcryptjs");

const run = async () => {
	try {
		await connectDB();

		const existing = await User.findOne({ email: "admin@physiolink.test" });
		if (existing) {
			console.log("Seed already applied");
			process.exit(0);
		}

		const password = "password123";
		const salt = await bcrypt.genSalt(10);
		const hashed = await bcrypt.hash(password, salt);

		const admin = new User({
			email: "admin@physiolink.test",
			password: hashed,
			role: "pt",
		});
		await admin.save();

		const profile = new Profile({
			user: admin._id,
			firstName: "Admin",
			lastName: "PT",
		});
		await profile.save();

		await User.findByIdAndUpdate(admin._id, { profile: profile._id });

		console.log(
			"Seed complete. Admin user created with email admin@physiolink.test and password",
			password,
		);
		process.exit(0);
	} catch (err) {
		console.error("Seed failed", err);
		process.exit(1);
	}
};

run();
