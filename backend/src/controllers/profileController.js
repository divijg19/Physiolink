// src/controllers/profileController.js
const profileService = require("../services/profileService");

// Get the profile for the currently logged-in user
exports.getCurrentUserProfile = async (req, res) => {
	try {
		const profile = await profileService.getCurrentUserProfile(req.user.id);
		res.json(profile);
	} catch (err) {
		if (err.status) return res.status(err.status).json({ msg: err.message });
		console.error(err);
		res.status(500).json({ msg: "Server Error" });
	}
};

// Create or update a user profile
exports.createOrUpdateProfile = async (req, res) => {
	try {
		const profile = await profileService.createOrUpdateProfile(
			req.user.id,
			req.body || {},
		);
		res.json(profile);
	} catch (err) {
		if (err.status) return res.status(err.status).json({ msg: err.message });
		console.error(err);
		res.status(500).json({ msg: "Server Error" });
	}
};
