# Back-end Coding Challenge

Welcome to this back-end Coding Challenge! We will ask you to complete the following exercise as part of our Interview process. It aims at assessing your coding skills proficiency, and your problem-solving mindset. It should take you no more than 2h.

You can submit this article as a public GitHub repository with the exercise and a short README explaining how to run it.  
You can write your code in **Go** (preferred), **JavaScript** or **Python**.

## We will assess
- If your solution works as expected
- The way you structure the code and its cleanliness
- Adherence to established standards when building API services
- The way you introduce code into your Git history

## Description
Let `User` and `Action` be two sets of records we store in our database. `User` contains basic information about our product’s users, and an `Action` is created when a user interacts with some features of our product, so we can keep an eye on our user’s behavior and improve our product accordingly.

### Data Structures
```json
User: {
  id: int,
  name: string,
  createdAt: date
}

Action: {
  id: int,
  type: string,
  userId: int,      // The ID of the User who performed this action
  targetUser: int,  // Supplied when "REFER_USER" action type
  createdAt: date
}
```

## Goal
Deliver a **very simple web server** capable of querying the relevant API endpoints described below.

### Additional Info
There is **no need to implement the database layer**, reading the file in memory at startup is sufficient.  
Use the two JSON files below containing a database sample for this exercise.

---

### 1) Fetch a User by ID
**Example response**
```json
{
  "id": "1234",
  "name": "John Doe",
  "createdAt": "2022-04-14T11:12:22.758Z"
}
```

---

### 2) Total number of actions of a User by ID
**Example response**
```json
{
  "count": 100
}
```

---

### 3) Next-action probability breakdown
Users perform different actions throughout the day when using our product. Right after doing an action A, users could perform action B, C, or A again.  
This endpoint returns the probability break-down of what actions our users typically do after doing action A, based on the entire Actions database.

**Example response**
```json
{
  "ADD_TO_CRM": 0.70,
  "REFER_USER": 0.20,
  "VIEW_CONVERSATION": 0.10
}
```

---

### 4) Referral Index
Users can refer other users (inviting them to use the product).  
When doing so, an activity with type `REFER_USER` is created by the existing user, and the ID of the new invited user is stored in the `targetUser` attribute.

We compute the **Referral Index** for a given user as the total number of individual users invited directly or indirectly by this user.  
Assume a user can be invited only **once**.

**Example response**
```json
{
  "1": 3,
  "2": 0,
  "3": 7
}
```
