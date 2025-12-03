// src/controllers/appointmentController.js
const appointmentService = require('../services/appointmentService');

// @desc    PT creates new available appointment slots
exports.createAvailability = async (req, res) => {
    try {
        const result = await appointmentService.createAvailability(req.user.id, req.body.slots || []);
        res.status(201).json(result);
    } catch (err) {
        if (err.status) return res.status(err.status).json({ msg: err.message });
        console.error(err);
        res.status(500).json({ msg: 'Server Error' });
    }
};

// @desc    Get all AVAILABLE appointment slots for a specific PT
exports.getTherapistAvailability = async (req, res) => {
    try {
        const availableSlots = await appointmentService.getTherapistAvailability(req.params.ptId);
        res.json(availableSlots);
    } catch (err) {
        console.error(err);
        res.status(500).json({ msg: 'Server Error' });
    }
};

// @desc    Patient books an available slot
exports.bookAppointment = async (req, res) => {
    try {
        const updated = await appointmentService.bookAppointment(req.params.id, req.user.id);
        return res.json(updated);
    } catch (err) {
        if (err.status) return res.status(err.status).json({ msg: err.message });
        console.error(err);
        res.status(500).json({ msg: 'Server Error' });
    }
};

// @desc PT updates appointment status (confirm/reject)
exports.updateAppointmentStatus = async (req, res) => {
    try {
        const appointment = await appointmentService.updateAppointmentStatus(req.params.id, req.user.id, req.body.status);
        return res.json(appointment);
    } catch (err) {
        if (err.status) return res.status(err.status).json({ msg: err.message });
        // If the service threw a reminder-creation wrapper error, log its cause
        if (err._cause) console.error('Reminder creation failed', err._cause);
        console.error(err);
        res.status(500).json({ msg: 'Server Error' });
    }
};

// @desc    Get appointments for the logged-in user
exports.getMySchedule = async (req, res) => {
    try {
        const appointments = await appointmentService.getMySchedule(req.user.id);
        res.json(appointments);
    } catch (err) {
        console.error(err);
        res.status(500).json({ msg: 'Server Error' });
    }
};
