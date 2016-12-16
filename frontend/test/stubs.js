export const adminUserStub = {
  id: 1,
  admin: true,
  email: 'hi@gnar.dog',
  name: 'Gnar Mike',
  username: 'gnardog',
};

export const scheduledQueryStub = {
  id: 1,
  interval: 60,
  name: 'Get all users',
  pack_id: 123,
  platform: 'darwin',
  query: 'SELECT * FROM users',
  query_id: 5,
  removed: false,
  snapshot: true,
};

export const userStub = {
  id: 1,
  admin: false,
  email: 'hi@gnar.dog',
  name: 'Gnar Mike',
  username: 'gnardog',
};

export const queryStub = {
  created_at: '2016-10-17T07:06:00Z',
  deleted: false,
  deleted_at: null,
  description: '',
  differential: false,
  id: 1,
  interval: 0,
  name: 'dev_query_1',
  platform: '',
  query: 'select * from processes',
  snapshot: false,
  updated_at: '2016-10-17T07:06:00Z',
  version: '',
};

export default { adminUserStub, queryStub, scheduledQueryStub, userStub };
