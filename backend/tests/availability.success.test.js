const request = require('supertest');
const mongoose = require('mongoose');
const { MongoMemoryServer } = require('mongodb-memory-server');

let mongoServer;
const app = require('../src/app');

const User = require('../src/models/User');
const Appointment = require('../src/models/Appointment');

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

const jwt = require('jsonwebtoken');
require('dotenv').config();
async function createPtAndToken() {
    const user = await User.create({ email: `${Math.random().toString(36).slice(2)}@pt`, password: 'pass', role: 'pt' });
    const payload = { user: { id: user._id.toString(), role: user.role } };
    const token = jwt.sign(payload, process.env.JWT_SECRET || 'testsecret', { expiresIn: '1h' });
    return { user, token };
}

test('createAvailability accepts non-overlapping slots', async () => {
    const { user: pt, token } = await createPtAndToken();

    const now = Date.now();
    const slot1 = { startTime: new Date(now + 3600 * 1000), endTime: new Date(now + 7200 * 1000) };
    const slot2 = { startTime: new Date(now + 7200 * 1000), endTime: new Date(now + 10800 * 1000) }; // adjacent but not overlapping

    const res = await request(app).post('/api/appointments/availability').set('x-auth-token', token).send({ slots: [slot1, slot2] });
    expect(res.status).toBe(201);

    const slots = await Appointment.find({ pt: pt._id });
    expect(slots.length).toBe(2);
});
