export const CONFIGURE_STAGED_QUERIES = 'CONFIGURE_STAGED_QUERIES';
export const STAGE_QUERY = 'STAGE_QUERY';
export const UNSTAGE_QUERY = 'UNSTAGE_QUERY';

export const unstageQuery = (query) => {
  return {
    type: UNSTAGE_QUERY,
    payload: { query },
  };
};
export const stageQuery = (query) => {
  return {
    type: STAGE_QUERY,
    payload: { query },
  };
};

export const configureStagedQueries = (configurationFormData) => {
  const {
    interval,
    logging_type: loggingType,
    platform,
    queries,
  } = configurationFormData;

  const configuration = {
    interval,
    logging_type: loggingType,
    platform,
    query_ids: queries.map(query => query.id),
  };

  return {
    type: CONFIGURE_STAGED_QUERIES,
    payload: { configuration },
  };
};

export default {
  configureStagedQueries,
  stageQuery,
  unstageQuery,
};
