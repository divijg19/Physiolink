// src/routes/appointments.js
const express = require("express");
const router = express.Router();
const auth = require("../../middleware/auth");
const isPt = require("../middleware/isPt");
const appointmentController = require("../controllers/appointmentController");

// @route   POST api/appointments/availability
// @desc    PT creates available slots
// @access  PT Only
router.post(
    "/availability",
    [auth, isPt],
    appointmentController.createAvailability,
);

// @route   GET api/appointments/me
// @desc    Get my schedule (for both PTs and Patients)
// @access  Private
router.get("/me", auth, appointmentController.getMySchedule);

// @route   GET api/appointments/availability/:ptId
// @desc    Get a specific PT's available slots
// @access  Private
router.get(
    "/availability/:ptId",
    auth,
    appointmentController.getTherapistAvailability,
);

// @route   PUT api/appointments/:id/book
// @desc    Patient books an appointment
// @access  Private
router.put("/:id/book", auth, appointmentController.bookAppointment);

// PT updates status of an appointment (confirm/reject)
router.put("/:id/status", [auth, isPt], appointmentController.updateAppointmentStatus);

module.exports = router;
