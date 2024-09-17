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

filter foreign value with local 
```sh
[
  {
    "$match": {
      "employeeId": {
        "$regex": "015",       // Partial match (substring) for employeeId
        "$options": "i"        // Case-insensitive match (optional)
      }
    }
  },
  {
    "$lookup": {
      "from": "meals",
      "localField": "_id",
      "foreignField": "consumerId",
      "as": "mealConsumption"
    }
  },
  {
    "$unwind": {
      "path": "$mealConsumption"
    }
  },
  {
    "$match": {
      "mealConsumption.month": 9,   // Replace with the desired month
      "mealConsumption.year": 2024  // Replace with the desired year
    }
  },
  {
    "$group": {
      "_id": "$_id",
      "totalMeals": {
        "$sum": "$mealConsumption.numberOfMeal"
      },
      "name": { "$first": "$name" },
      "employeeId": { "$first": "$employeeId" }
    }
  }
]

```