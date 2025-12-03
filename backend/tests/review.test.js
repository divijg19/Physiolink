const request = require("supertest");
const mongoose = require("mongoose");
const { MongoMemoryServer } = require("mongodb-memory-server");

let mongoServer;
const app = require("../src/app");

const User = require("../src/models/User");
const Profile = require("../src/models/Profile");
const Appointment = require("../src/models/Appointment");
const Review = require("../src/models/Review");

beforeAll(async () => {
	mongoServer = await MongoMemoryServer.create();
	const uri = mongoServer.getUri();
	await mongoose.disconnect();
	await mongoose.connect(uri);
});

afterAll(async () => {
	await mongoose.disconnect();
	await mongoServer.stop();
});

afterEach(async () => {
	await User.deleteMany({});
	await Profile.deleteMany({});
	await Appointment.deleteMany({});
	await Review.deleteMany({});
});

const jwt = require("jsonwebtoken");
require("dotenv").config();
async function createUserAndToken(role = "patient") {
	const user = await User.create({
		email: `${Math.random().toString(36).slice(2)}@test`,
		password: "pass",
		role,
	});
	// create profile for the user (mimic register flow)
	const profile = await Profile.create({
		user: user._id,
		firstName: "Test",
		lastName: "User",
	});
	user.profile = profile._id;
	await user.save();
	const payload = { user: { id: user._id.toString(), role: user.role } };
	const token = jwt.sign(payload, process.env.JWT_SECRET || "testsecret", {
		expiresIn: "1h",
	});
	return { user, token };
}

test("Review eligibility: cannot post without confirmed appointment, can post after confirm and rating updates", async () => {
	const { user: pt, token: ptToken } = await createUserAndToken("pt");
	const { token: patientToken } = await createUserAndToken("patient");

	// create availability
	const start = new Date(Date.now() + 48 * 3600 * 1000);
	const end = new Date(start.getTime() + 60 * 60 * 1000);
	const slot = await Appointment.create({
		pt: pt._id,
		startTime: start,
		endTime: end,
		status: "available",
	});

	// patient tries to post review before booking/confirm -> should be 403
	const badRes = await request(app)
		.post("/api/reviews")
		.set("x-auth-token", patientToken)
		.send({ therapistId: pt._id, rating: 5, comment: "Great!" });
	expect(badRes.status).toBe(403);

	// patient books
	await request(app)
		.put(`/api/appointments/${slot._id}/book`)
		.set("x-auth-token", patientToken)
		.expect(200);

	// PT confirms
	await request(app)
		.put(`/api/appointments/${slot._id}/status`)
		.set("x-auth-token", ptToken)
		.send({ status: "confirmed" })
		.expect(200);

	// now patient can submit a review
	const res = await request(app)
		.post("/api/reviews")
		.set("x-auth-token", patientToken)
		.send({ therapistId: pt._id, rating: 4, comment: "Helpful session" })
		.expect(201);
	expect(res.body).toHaveProperty("_id");

	// therapist's profile should be updated with avg rating
	const profile = await Profile.findOne({ user: pt._id });
	expect(profile.rating).toBeGreaterThanOrEqual(4);
});
