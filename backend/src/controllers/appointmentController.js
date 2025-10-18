// src/controllers/appointmentController.js
const Appointment = require("../models/Appointment");
const User = require("../models/User");

// @desc    PT creates new available appointment slots
exports.createAvailability = async (req, res) => {
    const { slots } = req.body; // Expect an array of { startTime, endTime }
    const ptId = req.user.id;

    try {
        const appointmentDocs = slots.map((slot) => ({
            pt: ptId,
            startTime: slot.startTime,
            endTime: slot.endTime,
            status: "available",
        }));

        await Appointment.insertMany(appointmentDocs);
        res.status(201).json({ msg: "Availability created successfully" });
    } catch (err) {
        // This will catch duplicate slot errors
        if (err.code === 11000) {
            return res
                .status(400)
                .json({ msg: "One or more of these time slots already exist." });
        }
        console.error(err);
        res.status(500).json({ msg: "Server Error" });
    }
};

// @desc    Get all AVAILABLE appointment slots for a specific PT
exports.getTherapistAvailability = async (req, res) => {
    try {
        const availableSlots = await Appointment.find({
            pt: req.params.ptId,
            status: "available",
        }).sort({ startTime: "asc" }); // Sort by time

        res.json(availableSlots);
    } catch (err) {
        console.error(err);
        res.status(500).json({ msg: "Server Error" });
    }
};

// @desc    Patient books an available slot
exports.bookAppointment = async (req, res) => {
    try {
        const appointment = await Appointment.findById(req.params.id);

        if (!appointment || appointment.status !== "available") {
            return res
                .status(404)
                .json({ msg: "This appointment slot is not available." });
        }

        appointment.patient = req.user.id;
        appointment.status = "booked";
        await appointment.save();

        res.json(appointment);
    } catch (err) {
        console.error(err);
        res.status(500).json({ msg: "Server Error" });
    }
};

// @desc PT updates appointment status (confirm/reject)
exports.updateAppointmentStatus = async (req, res) => {
    try {
        const { status } = req.body; // expected 'confirmed' or 'rejected'
        const appointment = await Appointment.findById(req.params.id);
        if (!appointment) return res.status(404).json({ msg: 'Appointment not found' });

        // Only the PT who owns the appointment can change its status
        if (appointment.pt.toString() !== req.user.id) {
            return res.status(403).json({ msg: 'Forbidden' });
        }

        if (!['confirmed', 'rejected'].includes(status)) {
            return res.status(400).json({ msg: 'Invalid status' });
        }

        appointment.status = status;
        await appointment.save();
        res.json(appointment);
    } catch (err) {
        console.error(err);
        res.status(500).json({ msg: 'Server Error' });
    }
};

// @desc    Get appointments for the logged-in user
exports.getMySchedule = async (req, res) => {
    try {
        const user = await User.findById(req.user.id);
        let query = {};

        if (user.role === "pt") {
            // THE FIX: PTs see ALL their slots, booked or available
            query = { pt: req.user.id };
        } else {
            // THE FIX: Patients ONLY see slots they have booked
            query = { patient: req.user.id };
        }

        const appointments = await Appointment.find(query)
            .populate({
                path: "patient",
                select: "profile",
                populate: { path: "profile" },
            })
            .populate({
                path: "pt",
                select: "profile",
                populate: { path: "profile" },
            })
            .sort({ startTime: "asc" });

        res.json(appointments);
    } catch (err) {
        console.error(err);
        res.status(500).json({ msg: "Server Error" });
    }
};
