const Review = require("../models/Review");
const Profile = require("../models/Profile");
const Appointment = require("../models/Appointment");

async function createReview(patientId, therapistId, rating, comment) {
	if (!therapistId || rating == null) {
		const err = new Error("therapistId and rating are required");
		err.status = 400;
		throw err;
	}

	const hadConfirmed = await Appointment.exists({
		pt: therapistId,
		patient: patientId,
		status: "confirmed",
	});
	if (!hadConfirmed) {
		const err = new Error(
			"Only patients with a confirmed appointment can leave a review",
		);
		err.status = 403;
		throw err;
	}

	const review = new Review({
		therapist: therapistId,
		patient: patientId,
		rating,
		comment,
	});
	await review.save();

	const agg = await Review.aggregate([
		{ $match: { therapist: review.therapist } },
		{ $group: { _id: "$therapist", avgRating: { $avg: "$rating" } } },
	]);
	const avgRating = agg[0] ? agg[0].avgRating : rating;

	await Profile.findOneAndUpdate(
		{ user: therapistId },
		{ rating: avgRating },
		{ upsert: false },
	);

	// Return populated review for API consumers
	const populated = await Review.findById(review._id)
		.populate({
			path: "therapist",
			select: "profile",
			populate: { path: "profile" },
		})
		.populate({
			path: "patient",
			select: "profile",
			populate: { path: "profile" },
		});
	return populated;
}

async function getReviewsForTherapist(therapistId) {
	const reviews = await Review.find({ therapist: therapistId }).populate(
		"patient",
		"profile",
	);
	return reviews;
}

module.exports = {
	createReview,
	getReviewsForTherapist,
};
