// src/controllers/therapistController.js
const therapistService = require('../services/therapistService');


// @desc    Get all verified physiotherapists
// Supports filtering (specialty, location), pagination (page, limit) and sorting (sort)
exports.getAllTherapists = async (req, res) => {
    try {
        const result = await therapistService.getAllTherapists(req.query || {});
        res.json(result);
    } catch (err) {
        console.error(err);
        res.status(500).json({ msg: 'Server Error' });
    }
};

// @desc Get a single therapist by id with availability
exports.getTherapistById = async (req, res) => {
    try {
        const therapist = await therapistService.getTherapistById(req.params.id, req.query.date);
        res.json(therapist);
    } catch (err) {
        if (err.status) return res.status(err.status).json({ msg: err.message });
        console.error(err);
        res.status(500).json({ msg: 'Server Error' });
    }
};
