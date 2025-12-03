const request = require("supertest");
const mongoose = require("mongoose");
const { MongoMemoryServer } = require("mongodb-memory-server");

let mongoServer;
const app = require("../src/app");

const User = require("../src/models/User");
const Profile = require("../src/models/Profile");
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
	await Profile.deleteMany({});
	await Appointment.deleteMany({});
});

const jwt = require("jsonwebtoken");
require("dotenv").config();
async function createUserAndToken(role = "patient") {
	const user = await User.create({
		email: `${Math.random().toString(36).slice(2)}@test`,
		password: "pass",
		role,
	});
	const profile = await Profile.create({
		user: user._id,
		firstName: "T",
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

test("GET /api/appointments/me returns PT full schedule and patient booked slots only", async () => {
	// create PT and patient
	const { user: pt, token: ptToken } = await createUserAndToken("pt");
	const { token: patientToken } = await createUserAndToken("patient");

	// create slots for PT: two available + one that will be booked
	const now = Date.now();
	const _slotA = await Appointment.create({
		pt: pt._id,
		startTime: new Date(now + 3600 * 1000),
		endTime: new Date(now + 7200 * 1000),
		status: "available",
	});
	const _slotB = await Appointment.create({
		pt: pt._id,
		startTime: new Date(now + 10800 * 1000),
		endTime: new Date(now + 14400 * 1000),
		status: "available",
	});
	const slotC = await Appointment.create({
		pt: pt._id,
		startTime: new Date(now + 18000 * 1000),
		endTime: new Date(now + 21600 * 1000),
		status: "available",
	});

	// patient books slotC via API
	await request(app)
		.put(`/api/appointments/${slotC._id}/book`)
		.set("x-auth-token", patientToken)
		.expect(200);

	// PT fetches their schedule -> should see all 3 slots
	const ptRes = await request(app)
		.get("/api/appointments/me")
		.set("x-auth-token", ptToken)
		.expect(200);
	expect(Array.isArray(ptRes.body)).toBeTruthy();
	expect(ptRes.body.length).toBe(3);

	// patient fetches their schedule -> should see only the booked slot (slotC)
	const patientRes = await request(app)
		.get("/api/appointments/me")
		.set("x-auth-token", patientToken)
		.expect(200);
	expect(Array.isArray(patientRes.body)).toBeTruthy();
	expect(patientRes.body.length).toBe(1);
	expect(patientRes.body[0]._id.toString()).toBe(slotC._id.toString());
});
