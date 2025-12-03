const Appointment = require('../models/Appointment');
const User = require('../models/User');

async function createAvailability(ptId, slots) {
    // Validate slots do not overlap with existing slots for this PT
    for (const slot of slots) {
        const start = new Date(slot.startTime);
        const end = new Date(slot.endTime);
        if (start >= end) {
            const err = new Error('Invalid slot times');
            err.status = 400;
            throw err;
        }
        const overlap = await Appointment.exists({
            pt: ptId,
            $or: [
                { $and: [{ startTime: { $lt: end } }, { endTime: { $gt: start } }] },
            ]
        });
        if (overlap) {
            const err = new Error('One or more slots overlap with existing availability');
            err.status = 400;
            throw err;
        }
    }

    const appointmentDocs = slots.map((slot) => ({
        pt: ptId,
        startTime: slot.startTime,
        endTime: slot.endTime,
        status: 'available',
    }));

    try {
        await Appointment.insertMany(appointmentDocs);
        return { msg: 'Availability created successfully' };
    } catch (err) {
        if (err.code === 11000) {
            const e = new Error('One or more of these time slots already exist.');
            e.status = 400;
            throw e;
        }
        throw err;
    }
}

async function getTherapistAvailability(ptId) {
    return Appointment.find({ pt: ptId, status: 'available' }).sort({ startTime: 'asc' });
}

async function bookAppointment(appointmentId, patientId) {
    // Atomic book
    const updated = await Appointment.findOneAndUpdate(
        { _id: appointmentId, status: 'available' },
        { $set: { patient: patientId, status: 'booked' } },
        { new: true }
    );
    if (!updated) {
        const exists = await Appointment.exists({ _id: appointmentId });
        const err = new Error(!exists ? 'Appointment not found' : 'This appointment slot is no longer available.');
        err.status = !exists ? 404 : 409;
        throw err;
    }
    // Return populated document for API responses
    const populated = await Appointment.findById(updated._id)
        .populate({ path: 'patient', select: 'profile', populate: { path: 'profile' } })
        .populate({ path: 'pt', select: 'profile', populate: { path: 'profile' } });
    return populated;
}

async function updateAppointmentStatus(appointmentId, userId, status) {
    const appointment = await Appointment.findById(appointmentId);
    if (!appointment) {
        const err = new Error('Appointment not found');
        err.status = 404;
        throw err;
    }
    if (appointment.pt.toString() !== userId) {
        const err = new Error('Forbidden');
        err.status = 403;
        throw err;
    }
    if (!['confirmed', 'rejected'].includes(status)) {
        const err = new Error('Invalid status');
        err.status = 400;
        throw err;
    }

    appointment.status = status;
    await appointment.save();

    // create reminder when confirmed
    if (status === 'confirmed') {
        try {
            const Reminder = require('../models/Reminder');
            const remindAt = new Date(appointment.startTime);
            remindAt.setHours(remindAt.getHours() - 24);
            await Reminder.create({
                appointment: appointment._id,
                patient: appointment.patient,
                pt: appointment.pt,
                message: `Reminder: appointment on ${appointment.startTime}`,
                remindAt,
            });
        } catch (err) {
            // Log but do not fail the status update if reminder creation fails
            console.error('Failed to create reminder', err);
        }
    }

    // Return populated appointment for responses
    const populated = await Appointment.findById(appointment._id)
        .populate({ path: 'patient', select: 'profile', populate: { path: 'profile' } })
        .populate({ path: 'pt', select: 'profile', populate: { path: 'profile' } });
    return populated;
}

async function getMySchedule(userId) {
    const user = await User.findById(userId);
    let query = {};
    if (user.role === 'pt') query = { pt: userId };
    else query = { patient: userId };

    return Appointment.find(query)
        .populate({ path: 'patient', select: 'profile', populate: { path: 'profile' } })
        .populate({ path: 'pt', select: 'profile', populate: { path: 'profile' } })
        .sort({ startTime: 'asc' });
}

module.exports = {
    createAvailability,
    getTherapistAvailability,
    bookAppointment,
    updateAppointmentStatus,
    getMySchedule,
};
