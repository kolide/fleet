import { combineReducers } from 'redux';

import hosts from './hosts/reducer';
import invites from './invites/reducer';
import labels from './labels/reducer';
import users from './users/reducer';

export default combineReducers({
  hosts,
  invites,
  labels,
  users,
});
