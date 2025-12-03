const request = require('supertest');
const path = require('path');
const mongoose = require('mongoose');
const { MongoMemoryServer } = require('mongodb-memory-server');
const SwaggerParser = require('@apidevtools/swagger-parser');
const Ajv = require('ajv');

let mongoServer;
const app = require('../src/app');

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

// Load and dereference OpenAPI once
let derefSpec;
let ajv;
beforeAll(async () => {
    const specPath = path.join(__dirname, '..', 'openapi.yaml');
    derefSpec = await SwaggerParser.dereference(specPath);
    ajv = new Ajv({ strict: false });
});

function getResponseSchema(pathTemplate, method) {
    const op = derefSpec.paths[pathTemplate]?.[method.toLowerCase()];
    if (!op) return null;
    const resp = op.responses && (op.responses['200'] || op.responses['201'] || op.responses['default']);
    const schema = resp && resp.content && resp.content['application/json'] && resp.content['application/json'].schema;
    return schema || null;
}

test('Validate /api/therapists list response against OpenAPI schema', async () => {
    // create minimal data: users and profiles
    const User = require('../src/models/User');
    const Profile = require('../src/models/Profile');
    const u1 = await User.create({ email: 'x1+schema@example.com', password: 'pass', role: 'pt' });
    const p1 = await Profile.create({ user: u1._id, firstName: 'X1', lastName: 'Schema' });
    u1.profile = p1._id; await u1.save();

    const patient = await User.create({ email: 'pat+schema@test', password: 'pass', role: 'patient' });
    const jwt = require('jsonwebtoken');
    const token = jwt.sign({ user: { id: patient._id.toString() } }, process.env.JWT_SECRET || 'testsecret');

    const res = await request(app).get('/api/therapists').set('x-auth-token', token).expect(200);
    const schema = getResponseSchema('/therapists', 'get');
    expect(schema).not.toBeNull();
    const validate = ajv.compile(schema);
    const ok = validate(res.body);
    if (!ok) console.error('Validation errors:', validate.errors);
    expect(ok).toBe(true);
});

test('Validate /api/therapists/{id} detail response against OpenAPI schema', async () => {
    const User = require('../src/models/User');
    const Profile = require('../src/models/Profile');
    const u1 = await User.create({ email: 'x2+schema@example.com', password: 'pass', role: 'pt' });
    const p1 = await Profile.create({ user: u1._id, firstName: 'X2', lastName: 'Schema' });
    u1.profile = p1._id; await u1.save();

    const patient = await User.create({ email: 'pat2+schema@test', password: 'pass', role: 'patient' });
    const jwt = require('jsonwebtoken');
    const token = jwt.sign({ user: { id: patient._id.toString() } }, process.env.JWT_SECRET || 'testsecret');

    // fetch list and take an id
    const listRes = await request(app).get('/api/therapists').set('x-auth-token', token).expect(200);
    const id = listRes.body.data[0]._id;
    const res = await request(app).get(`/api/therapists/${id}`).set('x-auth-token', token).expect(200);

    const schema = getResponseSchema('/therapists/{id}', 'get');
    expect(schema).not.toBeNull();
    const validate = ajv.compile(schema);
    const ok = validate(res.body);
    if (!ok) console.error('Detail validation errors:', validate.errors);
    expect(ok).toBe(true);
});

test('Validate POST /api/appointments/availability and GET /api/appointments/me', async () => {
    const User = require('../src/models/User');
    const Profile = require('../src/models/Profile');

    // create PT user
    const pt = await User.create({ email: 'pt+schema@example.com', password: 'pass', role: 'pt' });
    const ptProfile = await Profile.create({ user: pt._id, firstName: 'PT', lastName: 'Tester' });
    pt.profile = ptProfile._id; await pt.save();

    const jwt = require('jsonwebtoken');
    const ptToken = jwt.sign({ user: { id: pt._id.toString() } }, process.env.JWT_SECRET || 'testsecret');

    const start = new Date(Date.now() + 3600 * 1000).toISOString();
    const end = new Date(Date.now() + 7200 * 1000).toISOString();

    const payload = { slots: [{ startTime: start, endTime: end }] };

    const createRes = await request(app).post('/api/appointments/availability').set('x-auth-token', ptToken).send(payload).expect(201);
    const schemaPost = getResponseSchema('/appointments/availability', 'post');
    // spec may not declare a detailed response body for creation; only assert schema exists if present
    if (schemaPost) {
        const validate = ajv.compile(schemaPost);
        const ok = validate(createRes.body);
        if (!ok) console.error('Availability POST validation errors:', validate.errors);
        expect(ok).toBe(true);
    }

    // Now fetch /api/appointments/me as PT
    const meRes = await request(app).get('/api/appointments/me').set('x-auth-token', ptToken).expect(200);
    const schemaMe = getResponseSchema('/appointments/me', 'get');
    expect(schemaMe).not.toBeNull();
    const validateMe = ajv.compile(schemaMe);
    const okMe = validateMe(meRes.body);
    if (!okMe) console.error('Appointments/me validation errors:', validateMe.errors);
    expect(okMe).toBe(true);
});

test('Booking and review flow: book appointment, confirm, post review (schema checks when available)', async () => {
    const User = require('../src/models/User');
    const Profile = require('../src/models/Profile');

    // create PT and availability
    const pt = await User.create({ email: 'pt+flow@example.com', password: 'pass', role: 'pt' });
    const ptProfile = await Profile.create({ user: pt._id, firstName: 'PTflow', lastName: 'Tester' });
    pt.profile = ptProfile._id; await pt.save();

    const jwt = require('jsonwebtoken');
    const ptToken = jwt.sign({ user: { id: pt._id.toString() } }, process.env.JWT_SECRET || 'testsecret');

    const start = new Date(Date.now() + 3600 * 1000).toISOString();
    const end = new Date(Date.now() + 7200 * 1000).toISOString();
    const payload = { slots: [{ startTime: start, endTime: end }] };
    const createRes = await request(app).post('/api/appointments/availability').set('x-auth-token', ptToken).send(payload).expect(201);

    // find the slot via appointments/me
    const meRes = await request(app).get('/api/appointments/me').set('x-auth-token', ptToken).expect(200);
    const slot = meRes.body.find(s => s.startTime === start);
    expect(slot).toBeDefined();

    // create patient and book
    const patient = await User.create({ email: 'pat+flow@test', password: 'pass', role: 'patient' });
    const patientToken = jwt.sign({ user: { id: patient._id.toString() } }, process.env.JWT_SECRET || 'testsecret');

    const bookRes = await request(app).put(`/api/appointments/${slot._id}/book`).set('x-auth-token', patientToken).expect(200);
    // validate booking response schema if present
    const schemaBook = getResponseSchema('/appointments/{id}/book', 'put');
    if (schemaBook) {
        const validate = ajv.compile(schemaBook);
        const ok = validate(bookRes.body);
        if (!ok) console.error('Booking response validation errors:', validate.errors);
        expect(ok).toBe(true);
    }

    // PT confirms
    const confirmRes = await request(app).put(`/api/appointments/${slot._id}/status`).set('x-auth-token', ptToken).send({ status: 'confirmed' }).expect(200);
    const schemaStatus = getResponseSchema('/appointments/{id}/status', 'put');
    if (schemaStatus) {
        const validate = ajv.compile(schemaStatus);
        const ok = validate(confirmRes.body);
        if (!ok) console.error('Status response validation errors:', validate.errors);
        expect(ok).toBe(true);
    }

    // Patient posts a review
    const reviewPayload = { therapistId: pt._id.toString(), rating: 5, comment: 'Great session' };
    const reviewRes = await request(app).post('/api/reviews').set('x-auth-token', patientToken).send(reviewPayload).expect(201);
    const schemaReview = getResponseSchema('/reviews', 'post');
    if (schemaReview) {
        const validate = ajv.compile(schemaReview);
        const ok = validate(reviewRes.body);
        if (!ok) console.error('Review response validation errors:', validate.errors);
        expect(ok).toBe(true);
    }
});
