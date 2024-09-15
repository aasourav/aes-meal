filter from foreign field
```sh
db.getCollection('users').aggregate(
  [
    {
      $lookup: {
        from: 'meals',
        localField: '_id',
        foreignField: 'consumerId',
        as: 'mealsConsumed'
      }
    },
    {
      $addFields: {
        mealsConsumed: {
          $filter: {
            input: '$mealsConsumed',
            as: 'meal',
            cond: {
              $and: [
                {
                  $eq: ['$$meal.dayOfMonth', 15]
                },
                { $eq: ['$$meal.month', 9] },
                { $eq: ['$$meal.year', 2024] }
              ]
            }
          }
        }
      }
    },
    { $match: { mealsConsumed: { $ne: [] } } }
  ],
  { maxTimeMS: 60000, allowDiskUse: true }
);
```