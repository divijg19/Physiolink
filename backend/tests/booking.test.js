const request = require("supertest");
const mongoose = require("mongoose");
const { MongoMemoryServer } = require("mongodb-memory-server");

let mongoServer;
const app = require("../src/app");

const User = require("../src/models/User");
const Appointment = require("../src/models/Appointment");

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
	await Appointment.deleteMany({});
});

// helper to create a user and return a JWT
const jwt = require("jsonwebtoken");
require("dotenv").config();
async function createUserAndToken(role = "patient") {
	const user = await User.create({
		email: `${Math.random().toString(36).slice(2)}@test`,
		password: "pass",
		role,
	});
	const payload = { user: { id: user._id.toString() } };
	const token = jwt.sign(payload, process.env.JWT_SECRET || "testsecret", {
		expiresIn: "1h",
	});
	return { user, token };
}

test("concurrent booking: only one request should succeed", async () => {
	// create a PT and an availability slot
	const { user: ptUser } = await createUserAndToken("pt");

	// create an availability slot via direct model to keep test deterministic
	const slot = await Appointment.create({
		pt: ptUser._id,
		startTime: new Date(Date.now() + 3600 * 1000),
		endTime: new Date(Date.now() + 7200 * 1000),
		status: "available",
	});

	// create two patients and tokens
	const { token: t1 } = await createUserAndToken("patient");
	const { token: t2 } = await createUserAndToken("patient");

	// two concurrent booking attempts
	const p1 = request(app)
		.put(`/api/appointments/${slot._id}/book`)
		.set("x-auth-token", t1);
	const p2 = request(app)
		.put(`/api/appointments/${slot._id}/book`)
		.set("x-auth-token", t2);

	const results = await Promise.allSettled([p1, p2]);

	const statuses = results.map((r) => {
		if (r.status === "fulfilled") return r.value.status;
		return 500;
	});

	// There should be one 200 and one 409 (conflict) — order is nondeterministic
	expect(statuses.includes(200)).toBeTruthy();
	expect(statuses.includes(409)).toBeTruthy();
});
