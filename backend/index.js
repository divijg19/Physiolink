// index.js
const express = require("express");
const cors = require("cors");
require("dotenv").config();
const connectDB = require("./src/config/db");

// Initialize Express app
const app = express();

// Connect to Database
connectDB();

// Init Middleware
app.use(cors()); // Allows requests from our frontend
app.use(express.json()); // Allows us to accept JSON data in the body

// --- Define Routes ---
app.use("/api/auth", require("./src/routes/auth"));
app.use("/api/profile", require("./src/routes/profile"));
app.use("/api/therapists", require("./src/routes/therapists"));
app.use("/api/appointments", require("./src/routes/appointments"));

// A simple test route
app.get("/", (_req, res) => res.send("PhysioLink API is running!"));

// Centralized error handler (should be after all routes)
const errorHandler = require("./middleware/errorHandler");
app.use(errorHandler);

const PORT = process.env.PORT || 3001;

app.listen(PORT, () => console.log(`Server started on port ${PORT}`));
