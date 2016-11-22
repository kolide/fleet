import { combineReducers } from 'redux';

import ForgotPasswordPage from './ForgotPasswordPage/reducer';
import ManageHostsPage from './ManageHostsPage/reducer';
import PacksPages from './PacksPages/reducer';
import QueryPages from './QueryPages/reducer';
import ResetPasswordPage from './ResetPasswordPage/reducer';

export default combineReducers({
  ForgotPasswordPage,
  ManageHostsPage,
  PacksPages,
  QueryPages,
  ResetPasswordPage,
});

