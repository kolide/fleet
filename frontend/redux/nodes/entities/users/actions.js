import config from './config';
import Kolide from 'kolide';

export const REQUIRE_PASSWORD_RESET_FAILURE = "REQUIRE_PASSWORD_RESET_FAILURE";
export const REQUIRE_PASSWORD_RESET_REQUEST = "REQUIRE_PASSWORD_RESET_REQUEST";
export const REQUIRE_PASSWORD_RESET_SUCCESS = "REQUIRE_PASSWORD_RESET_SUCCESS";

export const requirePasswordResetFailure = (errors) => {
  return {
    type: REQUIRE_PASSWORD_RESET_FAILURE,
    payload: { errors },
  };
};

export const requirePasswordResetRequest = { type: REQUIRE_PASSWORD_RESET_REQUEST };

export const requirePasswordResetSuccess = { type: REQUIRE_PASSWORD_RESET_SUCCESS };

export const requirePasswordReset = (user, require = true) => {
  return (dispatch) => {
    dispatch(requirePasswordResetRequest);

    return Kolide.requirePasswordReset(user, require)
      .then((response) => {
        dispatch(requirePasswordResetSuccess);

        return response;
      })
      .catch((response) => {
        const { errors } = response;

        dispatch(requirePasswordResetFailure(errors));
        throw response;
      });
  };
};

export default {...config.actions, requirePasswordReset};
