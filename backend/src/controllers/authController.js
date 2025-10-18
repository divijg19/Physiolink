// src/controllers/authController.js
const User = require("../models/User");
const bcrypt = require("bcryptjs");
const jwt = require("jsonwebtoken");
const Profile = require("../models/Profile");

// The logic to register a new user
exports.registerUser = async (req, res) => {
	const { email, password, role } = req.body;

	try {
		// Check if user already exists
		let user = await User.findOne({ email });
		if (user) {
			return res.status(400).json({ msg: "User already exists" });
		}

		// Create a new user instance
		user = new User({ email, password, role });

		// Hash the password for security
		const salt = await bcrypt.genSalt(10);
		user.password = await bcrypt.hash(password, salt);

		// Save the user to the database
		await user.save();

		// Create an empty profile for the user (frontend will prompt to complete it)
		const profile = new Profile({
			user: user._id,
			firstName: "",
			lastName: "",
		});
		await profile.save();
		await User.findByIdAndUpdate(user._id, { profile: profile._id });

		// Create a JSON Web Token (JWT) to send back
		const payload = {
			user: {
				id: user.id,
				role: user.role,
			},
		};

		jwt.sign(
			payload,
			process.env.JWT_SECRET,
			{ expiresIn: "5h" },
			(err, token) => {
				if (err) throw err;
				res.json({ token, profile }); // Send the token and the newly created profile back to the user
			},
		);
	} catch (err) {
		console.error(err);
		res.status(500).json({ msg: "Server error" });
	}
};

// The logic to log in an existing user
exports.loginUser = async (req, res) => {
	const { email, password } = req.body;

	try {
		// 1. Check if a user with that email exists
		const user = await User.findOne({ email });
		if (!user) {
			// For security, give a generic error. Don't reveal if the user exists or not.
			return res.status(400).json({ msg: "Invalid Credentials" });
		}

		// 2. Compare the provided password with the hashed password in the database
		const isMatch = await bcrypt.compare(password, user.password);
		if (!isMatch) {
			return res.status(400).json({ msg: "Invalid Credentials" });
		}

		// 3. If credentials are correct, create and return a new token
		const payload = {
			user: {
				id: user.id,
				role: user.role,
			},
		};

		jwt.sign(
			payload,
			process.env.JWT_SECRET,
			{ expiresIn: "5h" }, // User stays logged in for 5 hours
			(err, token) => {
				if (err) throw err;
				res.json({ token });
			},
		);
	} catch (err) {
		console.error(err);
		res.status(500).json({ msg: "Server error" });
	}
};
