// src/controllers/therapistController.js
const User = require("../models/User");

// @desc    Get all verified physiotherapists
exports.getAllTherapists = async (req, res) => {
    try {
        // Accept optional query parameters for filtering
        const { specialty, location } = req.query;
        const userQuery = { role: 'pt' };

        // We'll do filtering on the profile fields using aggregation or a two-step query
        let therapists = await User.find(userQuery).select('-password').populate('profile');

        if (specialty) {
            therapists = therapists.filter(t => t.profile && t.profile.specialty && t.profile.specialty.toLowerCase().includes(specialty.toLowerCase()));
        }
        if (location) {
            therapists = therapists.filter(t => t.profile && t.profile.location && t.profile.location.toLowerCase().includes(location.toLowerCase()));
        }

        res.json(therapists);
    } catch (err) {
        console.error(err);
        res.status(500).json({ msg: 'Server Error' });
    }
};
