const request = require('supertest');
const mongoose = require('mongoose');
const { MongoMemoryServer } = require('mongodb-memory-server');

let mongoServer;
const app = require('../src/app');

const User = require('../src/models/User');
const Profile = require('../src/models/Profile');

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
});

test('Contract smoke: health, therapists list and therapist detail', async () => {
    // health check
    await request(app).get('/').expect(200);

    // create two therapists
    const u1 = await User.create({ email: 'alice+ct@example.com', password: 'pass', role: 'pt' });
    const p1 = await Profile.create({ user: u1._id, firstName: 'Alice', lastName: 'Smith', specialty: 'Orthopedics', location: 'NY' });
    u1.profile = p1._id; await u1.save();

    const u2 = await User.create({ email: 'bob+ct@example.com', password: 'pass', role: 'pt' });
    const p2 = await Profile.create({ user: u2._id, firstName: 'Bob', lastName: 'Jones', specialty: 'Sports', location: 'CA' });
    u2.profile = p2._id; await u2.save();

    // create a patient and token
    const jwt = require('jsonwebtoken');
    const patient = await User.create({ email: 'pat+ct@test', password: 'pass', role: 'patient' });
    const token = jwt.sign({ user: { id: patient._id.toString() } }, process.env.JWT_SECRET || 'testsecret');

    // GET /api/therapists
    const res = await request(app).get('/api/therapists').set('x-auth-token', token).expect(200);
    expect(res.body).toHaveProperty('data');
    expect(Array.isArray(res.body.data)).toBe(true);
    expect(res.body.data.length).toBeGreaterThanOrEqual(2);

    // GET /api/therapists/:id (detail)
    const first = res.body.data[0];
    const detailRes = await request(app).get(`/api/therapists/${first._id}`).set('x-auth-token', token).expect(200);
    expect(detailRes.body).toHaveProperty('profile');
    expect(detailRes.body.profile).toHaveProperty('firstName');
});
