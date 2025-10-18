// src/routes/profile.js
const express = require("express");
const router = express.Router();
const auth = require("../../middleware/auth"); // Import our auth middleware
const profileController = require("../controllers/profileController");

// @route   GET api/profile/me
// @desc    Get current user's profile
// @access  Private (notice the 'auth' middleware is passed in)
router.get("/me", auth, profileController.getCurrentUserProfile);

// @route   POST api/profile
// @desc    Create or update user profile
// @access  Private
router.post("/", auth, profileController.createOrUpdateProfile);

module.exports = router;
