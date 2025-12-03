const Profile = require("../models/Profile");
const User = require("../models/User");

async function getCurrentUserProfile(userId) {
	const profile = await Profile.findOne({ user: userId }).populate("user", [
		"email",
		"role",
	]);
	if (!profile) {
		const err = new Error("There is no profile for this user");
		err.status = 400;
		throw err;
	}
	return profile;
}

async function createOrUpdateProfile(userId, body) {
	const {
		firstName,
		lastName,
		age,
		gender,
		condition,
		goals,
		specialty,
		bio,
		credentials,
		location,
		profileImageUrl,
	} = body;

	const profileFields = { user: userId, firstName, lastName };
	if (age) profileFields.age = age;
	if (gender) profileFields.gender = gender;
	if (condition) profileFields.condition = condition;
	if (goals) profileFields.goals = goals;
	if (specialty) profileFields.specialty = specialty;
	if (bio) profileFields.bio = bio;
	if (credentials) profileFields.credentials = credentials;
	if (location) profileFields.location = location;
	if (profileImageUrl) profileFields.profileImageUrl = profileImageUrl;

	let profile = await Profile.findOne({ user: userId });
	if (profile) {
		profile = await Profile.findOneAndUpdate(
			{ user: userId },
			{ $set: profileFields },
			{ new: true },
		);
	} else {
		profile = new Profile(profileFields);
		await profile.save();
	}

	await User.findByIdAndUpdate(userId, { profile: profile._id });
	return profile;
}

module.exports = {
	getCurrentUserProfile,
	createOrUpdateProfile,
};
