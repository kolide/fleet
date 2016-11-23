import { pull } from 'lodash';

import { CONFIGURE_STAGED_QUERIES, STAGE_QUERY, UNSTAGE_QUERY } from './actions';

export const initialState = {
  stagedQueries: [],
  configuredQueryIDs: [],
  configurations: [],
};

const queryIDs = (queries) => {
  return queries.map((query) => query.id);
}

const reducer = (state = initialState, { type, payload }) => {
  switch (type) {
    case CONFIGURE_STAGED_QUERIES:
      return {
        ...state,
        stagedQueries: [],
        configuredQueryIDs: [
          ...state.configuredQueryIDs,
          ...queryIDs(state.stagedQueries),
        ],
        configurations: [
          ...state.configurations,
          payload.configuration,
        ],
      };
    case STAGE_QUERY:
      return {
        ...state,
        stagedQueries: [
          ...state.stagedQueries,
          payload.query,
        ],
      };
    case UNSTAGE_QUERY:
      return {
        ...state,
        stagedQueries: [
          ...pull(state.stagedQueries, payload.query),
        ],
      };
    default:
      return state;
  }
};

export default reducer;
