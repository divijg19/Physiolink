const express = require("express");
const router = express.Router();
const auth = require("../../middleware/auth");
const reminderController = require("../controllers/reminderController");

router.get("/me", auth, reminderController.getMyReminders);

module.exports = router;
