export const adminUserStub = {
  id: 1,
  admin: true,
  email: 'hi@gnar.dog',
  name: 'Gnar Mike',
  username: 'gnardog',
};

export const scheduledQueryStub = {
  id: 1,
  query_id: 5,
  pack_id: 123,
  interval: 60,
  snapshot: true,
};

export const userStub = {
  id: 1,
  admin: false,
  email: 'hi@gnar.dog',
  name: 'Gnar Mike',
  username: 'gnardog',
};

export default { adminUserStub, scheduledQueryStub, userStub };
