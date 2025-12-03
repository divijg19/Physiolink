const mongoose = require('mongoose');

const ReviewSchema = new mongoose.Schema({
    therapist: { type: mongoose.Schema.Types.ObjectId, ref: 'user', required: true },
    patient: { type: mongoose.Schema.Types.ObjectId, ref: 'user', required: true },
    rating: { type: Number, required: true, min: 0, max: 5 },
    comment: { type: String },
    date: { type: Date, default: Date.now },
});

// index to improve lookup performance when aggregating/retrieving reviews per therapist
ReviewSchema.index({ therapist: 1 });

module.exports = mongoose.model('review', ReviewSchema);
