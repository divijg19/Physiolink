const mongoose = require('mongoose');

async function getAllTherapists(params = {}) {
    const { specialty, location, page = 1, limit = 20, sort, date, available } = params;

    const pageNum = Math.max(1, parseInt(page, 10) || 1);
    const limitNum = Math.max(1, Math.min(100, parseInt(limit, 10) || 20));
    const skip = (pageNum - 1) * limitNum;

    const basePipeline = [
        { $match: { role: 'pt' } },
        {
            $lookup: {
                from: 'profiles',
                localField: 'profile',
                foreignField: '_id',
                as: 'profile',
            },
        },
        { $unwind: { path: '$profile', preserveNullAndEmptyArrays: true } },
    ];

    const filters = [];
    if (specialty) filters.push({ 'profile.specialty': { $regex: specialty, $options: 'i' } });
    if (location) filters.push({ 'profile.location': { $regex: location, $options: 'i' } });
    if (filters.length) basePipeline.push({ $match: { $and: filters } });

    const sortStage = {};
    if (sort) {
        const sortFields = Array.isArray(sort) ? sort : String(sort).split(',');
        sortFields.forEach(field => {
            let dir = 1;
            let fname = field;
            if (field.startsWith('-')) {
                dir = -1;
                fname = field.slice(1);
            }
            sortStage[fname] = dir;
        });
    } else {
        sortStage.date = -1;
    }

    const dateFilter = date;

    basePipeline.push({
        $lookup: {
            from: 'appointments',
            let: { userId: '$_id' },
            pipeline: [
                { $match: { $expr: { $and: [{ $eq: ['$pt', '$$userId'] }, { $eq: ['$status', 'available'] }] } } },
                ...(dateFilter ? [{ $addFields: { startDateString: { $dateToString: { format: '%Y-%m-%d', date: '$startTime' } } } }, { $match: { startDateString: dateFilter } }] : []),
                { $sort: { startTime: 1 } },
                { $limit: 20 }
            ],
            as: 'availableSlots'
        }
    });

    const requireAvailable = available === 'true' || available === true;

    basePipeline.push({ $addFields: { availableSlotsCount: { $size: '$availableSlots' } } });
    basePipeline.push({
        $lookup: {
            from: 'reviews',
            let: { userId: '$_id' },
            pipeline: [
                { $match: { $expr: { $eq: ['$therapist', '$$userId'] } } },
                { $count: 'count' }
            ],
            as: 'reviewAgg'
        }
    });
    basePipeline.push({ $addFields: { reviewCount: { $ifNull: [{ $arrayElemAt: ['$reviewAgg.count', 0] }, 0] } } });
    basePipeline.push({ $project: { password: 0, reviewAgg: 0 } });

    if (requireAvailable) basePipeline.push({ $match: { availableSlotsCount: { $gt: 0 } } });

    const pipeline = [
        ...basePipeline,
        {
            $facet: {
                results: [
                    { $sort: sortStage },
                    { $skip: skip },
                    { $limit: limitNum },
                ],
                totalCount: [
                    { $count: 'count' }
                ]
            }
        }
    ];

    const aggResult = await mongoose.model('user').aggregate(pipeline);
    const results = (aggResult[0] && aggResult[0].results) || [];
    const totalCount = (aggResult[0] && aggResult[0].totalCount && aggResult[0].totalCount[0]) ? aggResult[0].totalCount[0].count : 0;
    const totalPages = Math.ceil(totalCount / limitNum) || 1;

    return { data: results, total: totalCount, page: pageNum, totalPages };
}

async function getTherapistById(therapistId, dateFilter) {
    const pipeline = [
        { $match: { _id: new mongoose.Types.ObjectId(therapistId), role: 'pt' } },
        {
            $lookup: {
                from: 'profiles',
                localField: 'profile',
                foreignField: '_id',
                as: 'profile'
            }
        },
        { $unwind: { path: '$profile', preserveNullAndEmptyArrays: true } },
        {
            $lookup: {
                from: 'appointments',
                let: { userId: '$_id' },
                pipeline: [
                    { $match: { $expr: { $and: [{ $eq: ['$pt', '$$userId'] }, { $eq: ['$status', 'available'] }] } } },
                    ...(dateFilter ? [{ $addFields: { startDateString: { $dateToString: { format: '%Y-%m-%d', date: '$startTime' } } } }, { $match: { startDateString: dateFilter } }] : []),
                    { $sort: { startTime: 1 } }
                ],
                as: 'availableSlots'
            }
        },
        {
            $lookup: {
                from: 'reviews',
                let: { userId: '$_id' },
                pipeline: [
                    { $match: { $expr: { $eq: ['$therapist', '$$userId'] } } },
                    { $count: 'count' }
                ],
                as: 'reviewAgg'
            }
        },
        { $addFields: { reviewCount: { $ifNull: [{ $arrayElemAt: ['$reviewAgg.count', 0] }, 0] } } },
        { $project: { password: 0, reviewAgg: 0 } },
    ];

    const result = await mongoose.model('user').aggregate(pipeline);
    if (!result || !result[0]) {
        const err = new Error('Therapist not found');
        err.status = 404;
        throw err;
    }
    return result[0];
}

module.exports = {
    getAllTherapists,
    getTherapistById,
};
