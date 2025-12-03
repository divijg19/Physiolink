const Reminder = require("../models/Reminder");

// GET /api/reminders/me
exports.getMyReminders = async (req, res) => {
	try {
		const reminders = await Reminder.find({ patient: req.user.id }).populate(
			"appointment",
		);
		res.json(reminders);
	} catch (err) {
		console.error(err);
		res.status(500).json({ msg: "Server Error" });
	}
};
