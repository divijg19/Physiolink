const express = require("express");
const cors = require("cors");
const connectDB = require("./config/db");

// Create express app
const app = express();

// Connect to DB (will use real DB in dev). Tests use an in-memory server and manage connection themselves.
if (process.env.NODE_ENV !== "test") {
	connectDB();
}

app.use(cors());
app.use(express.json());

// Routes (resolve relative to this file to avoid resolver issues in test environments)
const path = require("node:path");
app.use("/api/auth", require(path.join(__dirname, "routes", "auth")));
app.use("/api/profile", require(path.join(__dirname, "routes", "profile")));
app.use(
	"/api/therapists",
	require(path.join(__dirname, "routes", "therapists")),
);
app.use(
	"/api/appointments",
	require(path.join(__dirname, "routes", "appointments")),
);
app.use("/api/reviews", require(path.join(__dirname, "routes", "reviews")));
app.use("/api/reminders", require(path.join(__dirname, "routes", "reminders")));

// Health check
app.get("/", (_req, res) => res.send("PhysioLink API is running!"));

const errorHandler = require("../middleware/errorHandler");
app.use(errorHandler);

module.exports = app;
