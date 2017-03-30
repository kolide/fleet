import Kolide from 'kolide';
import config from 'redux/nodes/entities/config_options/config';


const { actions } = config;

export const RESET_OPTIONS_START = 'RESET_OPTIONS_START';
export const RESET_OPTIONS_SUCCESS = 'RESET_OPTIONS_SUCCESS';
export const RESET_OPTIONS_FAILURE = 'RESET_OPTIONS_FAILURE';

export const resetOptionsStart = { type: RESET_OPTIONS_START };
export const resetOptionsSuccess = (config_options) => {
  return { type: RESET_OPTIONS_SUCCESS, payload:  {config_options} };
};
export const resetOptionsFailure = (errors) => {
  return { type: RESET_OPTIONS_FAILURE, payload: { errors } };
};

export const resetOptions = () => {

  return (dispatch) => {
    dispatch(resetOptionsStart);
    return Kolide.configOptions.reset()
      .then((opts) => {
        return dispatch(resetOptionsSuccess(opts));
      })
      .catch((error) => {
        const formattedErrors = formatApiErrors(error);
        dispatch(resetOptionsFailure(formattedErrors));
        throw formattedErrors;
      });
  };
};

export default {
  ...actions,
  resetOptions,
};
