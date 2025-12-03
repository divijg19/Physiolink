// src/routes/therapists.js
const express = require("express");
const router = express.Router();
const auth = require("../../middleware/auth");
const therapistController = require("../controllers/therapistController");

// @route   GET api/therapists
// @desc    Get all therapists
// @access  Private (must be logged in to see therapists)
router.get("/", auth, therapistController.getAllTherapists);
// GET api/therapists/:id
router.get('/:id', auth, therapistController.getTherapistById);

module.exports = router;
