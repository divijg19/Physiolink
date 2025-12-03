const mongoose = require('mongoose');

const ReminderSchema = new mongoose.Schema({
    appointment: { type: mongoose.Schema.Types.ObjectId, ref: 'appointment', required: true },
    patient: { type: mongoose.Schema.Types.ObjectId, ref: 'user', required: true },
    pt: { type: mongoose.Schema.Types.ObjectId, ref: 'user', required: true },
    message: { type: String },
    remindAt: { type: Date, required: true },
    sent: { type: Boolean, default: false },
    dateCreated: { type: Date, default: Date.now },
});

module.exports = mongoose.model('reminder', ReminderSchema);
