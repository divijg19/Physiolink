// middleware/auth.js
const jwt = require("jsonwebtoken");
require("dotenv").config();

module.exports = (req, res, next) => {
	// 1. Get the token from the header
	const token = req.header("x-auth-token");

	// 2. Check if there is no token
	if (!token) {
		return res.status(401).json({ msg: "No token, authorization denied" });
	}

	// 3. Verify the token
	try {
		const decoded = jwt.verify(token, process.env.JWT_SECRET);
		req.user = decoded.user; // Add the user payload to the request object
		next(); // Move on to the next piece of middleware or the route handler
	} catch (_err) {
		res.status(401).json({ msg: "Token is not valid" });
	}
};
