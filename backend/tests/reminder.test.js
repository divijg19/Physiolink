const request = require("supertest");
const mongoose = require("mongoose");
const { MongoMemoryServer } = require("mongodb-memory-server");

let mongoServer;
const app = require("../src/app");

const User = require("../src/models/User");
const Appointment = require("../src/models/Appointment");
const Reminder = require("../src/models/Reminder");

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
	await Reminder.deleteMany({});
});

const jwt = require("jsonwebtoken");
require("dotenv").config();
async function createUserAndToken(role = "patient") {
	const user = await User.create({
		email: `${Math.random().toString(36).slice(2)}@test`,
		password: "pass",
		role,
	});
	const payload = { user: { id: user._id.toString(), role: user.role } };
	const token = jwt.sign(payload, process.env.JWT_SECRET || "testsecret", {
		expiresIn: "1h",
	});
	return { user, token };
}

test("PT confirming an appointment creates a Reminder 24 hours before", async () => {
	// create PT and patient
	const { user: pt, token: ptToken } = await createUserAndToken("pt");
	const { user: patient, token: patientToken } =
		await createUserAndToken("patient");

	// create availability slot
	const start = new Date(Date.now() + 48 * 3600 * 1000); // 48 hours from now
	const end = new Date(start.getTime() + 60 * 60 * 1000);
	const slot = await Appointment.create({
		pt: pt._id,
		startTime: start,
		endTime: end,
		status: "available",
	});

	// patient books the slot
	await request(app)
		.put(`/api/appointments/${slot._id}/book`)
		.set("x-auth-token", patientToken)
		.expect(200);

	// PT confirms the appointment
	const confirmRes = await request(app)
		.put(`/api/appointments/${slot._id}/status`)
		.set("x-auth-token", ptToken)
		.send({ status: "confirmed" })
		.expect(200);
	expect(confirmRes.body.status).toBe("confirmed");

	// verify reminder created
	const reminder = await Reminder.findOne({ appointment: slot._id });
	expect(reminder).not.toBeNull();
	expect(reminder.patient.toString()).toBe(patient._id.toString());
	expect(reminder.pt.toString()).toBe(pt._id.toString());

	const expectedRemindAt = new Date(start);
	expectedRemindAt.setHours(expectedRemindAt.getHours() - 24);
	// allow small tolerance (2s)
	const diff = Math.abs(
		reminder.remindAt.getTime() - expectedRemindAt.getTime(),
	);
	expect(diff).toBeLessThan(2000);
});
