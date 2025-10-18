// src/controllers/profileController.js
const Profile = require("../models/Profile");
const User = require("../models/User");

// Get the profile for the currently logged-in user
exports.getCurrentUserProfile = async (req, res) => {
	try {
		const profile = await Profile.findOne({ user: req.user.id }).populate(
			"user",
			["email", "role"],
		);
		if (!profile) {
			return res.status(400).json({ msg: "There is no profile for this user" });
		}
		res.json(profile);
	} catch (err) {
		console.error(err);
		res.status(500).json({ msg: "Server Error" });
	}
};

// Create or update a user profile
exports.createOrUpdateProfile = async (req, res) => {
	// Destructure all possible fields from the request body
	const {
		firstName,
		lastName,
		age,
		gender,
		condition,
		goals,
		specialty,
		bio,
		credentials,
	} = req.body;
	const profileFields = { user: req.user.id, firstName, lastName };

	if (age) profileFields.age = age;
	if (gender) profileFields.gender = gender;
	if (condition) profileFields.condition = condition;
	if (goals) profileFields.goals = goals;
	if (specialty) profileFields.specialty = specialty;
	if (bio) profileFields.bio = bio;
	if (credentials) profileFields.credentials = credentials;

	try {
		let profile = await Profile.findOne({ user: req.user.id });
		if (profile) {
			// Update existing profile
			profile = await Profile.findOneAndUpdate(
				{ user: req.user.id },
				{ $set: profileFields },
				{ new: true },
			);
		} else {
			// Create new profile
			profile = new Profile(profileFields);
			await profile.save();
		}

		// Link the profile back to the User document
		await User.findByIdAndUpdate(req.user.id, { profile: profile._id });

		res.json(profile);
	} catch (err) {
		console.error(err);
		res.status(500).json({ msg: "Server Error" });
	}
};
