import { noop } from 'lodash';
import { normalize, arrayOf } from 'normalizr';

const initialState = {
  loading: false,
  errors: {},
  entities: {},
};

const reduxConfig = ({
  entityName,
  loadFunc,
  parseFunc = noop,
  schema,
}) => {
  const actionTypes = {
    LOAD_FAILURE: `${entityName}_LOAD_FAILURE`,
    LOAD_REQUEST: `${entityName}_LOAD_REQUEST`,
    LOAD_SUCCESS: `${entityName}_LOAD_SUCCESS`,
  };

  const loadFailure = (errors) => {
    return {
      type: actionTypes.LOAD_FAILURE,
      payload: { errors },
    };
  };
  const loadRequest = { type: actionTypes.LOAD_REQUEST };
  const loadSuccess = (entities) => {
    return {
      type: actionTypes.LOAD_SUCCESS,
      payload: { entities },
    };
  };

  const parsedResponse = (responseArray) => {
    return responseArray.map(response => {
      return parseFunc(response);
    });
  };

  const load = (...args) => {
    return (dispatch) => {
      dispatch(loadRequest);

      return loadFunc(...args)
        .then(response => {
          if (!response) return [];

          const { entities } = normalize(parsedResponse(response), arrayOf(schema));

          return dispatch(loadSuccess(entities));
        })
        .catch(response => {
          const { errors } = response;

          dispatch(loadFailure(errors));
          throw response;
        });
    };
  };

  const actions = {
    load,
  };

  const reducer = (state = initialState, { type, payload }) => {
    switch (type) {
      case actionTypes.LOAD_REQUEST:
        return {
          ...state,
          loading: true,
        };
      case actionTypes.LOAD_SUCCESS:
        return {
          ...state,
          loading: false,
          entities: {
            ...state.entities,
            ...payload.entities[entityName],
          },
        };
      case actionTypes.LOAD_FAILURE:
        return {
          ...state,
          loading: false,
          errors: {
            ...payload.errors,
          },
        };
      default:
        return state;
    }
  };

  return {
    actions,
    reducer,
  };
};

export default reduxConfig;
