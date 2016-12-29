import {
  REQUIRE_PASSWORD_RESET_FAILURE,
  REQUIRE_PASSWORD_RESET_REQUEST,
  REQUIRE_PASSWORD_RESET_SUCCESS
} from './actions';
import config, { initialState } from './config';

var reducer = (state = initialState, {type, payload}) => {
  switch (type) {
  case REQUIRE_PASSWORD_RESET_REQUEST:
    return {
      ...state,
      loading: true,
    };
  case REQUIRE_PASSWORD_RESET_SUCCESS:
    return {
      ...state,
      loading: false,
    };
  case REQUIRE_PASSWORD_RESET_FAILURE:
    return {
      ...state,
      loading:false,
      errors: payload.errors,
    };
  default:
    return config.reducer(state, {type, payload})
  }
}

export default reducer;
