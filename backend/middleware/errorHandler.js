// middleware/errorHandler.js
// Centralized Express error handler — returns JSON and logs detailed stack in development
module.exports = (err, _req, res, _next) => {
	console.error("Unhandled Error:", err);
	const response = { msg: "Server Error" };
	if (process.env.NODE_ENV === "development") {
		response.error = err.message;
		response.stack = err.stack;
	}
	res.status(500).json(response);
};
