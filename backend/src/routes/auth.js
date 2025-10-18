// src/routes/auth.js
const express = require("express");
const router = express.Router();
const authController = require("../controllers/authController");

// @route   POST api/auth/register
// @desc    Register a user
// @access  Public
router.post("/register", authController.registerUser);

// @route   POST api/auth/login
// @desc    Authenticate user & get token
// @access  Public
router.post("/login", authController.loginUser); // <-- ADD THIS LINE

module.exports = router;
