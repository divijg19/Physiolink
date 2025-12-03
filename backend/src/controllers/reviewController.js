const reviewService = require("../services/reviewService");

// POST /api/reviews
exports.createReview = async (req, res) => {
	try {
		const { therapistId, rating, comment } = req.body;
		const patientId = req.user.id;
		const review = await reviewService.createReview(
			patientId,
			therapistId,
			rating,
			comment,
		);
		res.status(201).json(review);
	} catch (err) {
		if (err.status) return res.status(err.status).json({ msg: err.message });
		console.error(err);
		res.status(500).json({ msg: "Server Error" });
	}
};

// GET /api/reviews/:therapistId
exports.getReviewsForTherapist = async (req, res) => {
	try {
		const therapistId = req.params.therapistId;
		const reviews = await reviewService.getReviewsForTherapist(therapistId);
		res.json(reviews);
	} catch (err) {
		console.error(err);
		res.status(500).json({ msg: "Server Error" });
	}
};
