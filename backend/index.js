// index.js
const express = require("express");
const cors = require("cors");
require("dotenv").config();
const connectDB = require("./src/config/db");

// Use the app factory so tests can import it without starting a server
const app = require('./src/app');

const PORT = process.env.PORT || 3001;

app.listen(PORT, () => console.log(`Server started on port ${PORT}`));
