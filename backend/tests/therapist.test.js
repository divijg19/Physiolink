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
    // clean up
    await User.deleteMany({});
    await Profile.deleteMany({});
});

test('GET /api/therapists returns paginated therapists and supports filtering', async () => {
    // create two users and profiles
    const u1 = await User.create({ email: 'alice@example.com', password: 'pass', role: 'pt' });
    const p1 = await Profile.create({ user: u1._id, firstName: 'Alice', lastName: 'Smith', specialty: 'Orthopedics', location: 'NY' });
    u1.profile = p1._id; await u1.save();

    const u2 = await User.create({ email: 'bob@example.com', password: 'pass', role: 'pt' });
    const p2 = await Profile.create({ user: u2._id, firstName: 'Bob', lastName: 'Jones', specialty: 'Sports', location: 'CA' });
    u2.profile = p2._id; await u2.save();

    // create a patient and token to call protected endpoint
    const jwt = require('jsonwebtoken');
    const patient = await User.create({ email: 'pat@test', password: 'pass', role: 'patient' });
    const token = jwt.sign({ user: { id: patient._id.toString() } }, process.env.JWT_SECRET || 'testsecret');
    const resAll = await request(app).get('/api/therapists').set('x-auth-token', token).expect(200);
    expect(resAll.body).toHaveProperty('data');
    expect(resAll.body.data.length).toBe(2);

    // filter by specialty
    const resFilter = await request(app).get('/api/therapists').set('x-auth-token', token).query({ specialty: 'sports' }).expect(200);
    expect(resFilter.body.data.length).toBe(1);
    expect(resFilter.body.data[0].profile.specialty.toLowerCase()).toContain('sports');
});
