// src/middleware/isPt.js
const User = require("../models/User");

module.exports = async (req, res, next) => {
	// We assume the 'auth' middleware has already run and added req.user
	if (!req.user) {
		return res.status(401).json({ msg: "Authorization denied" });
	}

	try {
		const user = await User.findById(req.user.id);
		if (user.role !== "pt") {
			return res
				.status(403)
				.json({ msg: "Access forbidden: User is not a Physiotherapist" });
		}
		next();
	} catch (err) {
		console.error(err);
		res.status(500).json({ msg: "Server Error" });
	}
};
