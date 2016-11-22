import expect from 'expect';

import reducer, { initialState } from './reducer';
import * as actions from './actions';

const query = {
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

describe('PacksPages - reducer', () => {
  it('sets the initial state', () => {
    expect(reducer(undefined, { type: 'SOME ACTION' })).toEqual(initialState);
  });

  it('stages queries', () => {
    const stageQueryAction = actions.stageQuery(query);

    expect(reducer(initialState, stageQueryAction)).toEqual({
      ...initialState,
      stagedQueries: [query],
    });
  });

  it('unstages queries', () => {
    const stagedQueryState = {
      ...initialState,
      stagedQueries: [query],
    };
    const unstageQueryAction = actions.unstageQuery(query);

    expect(reducer(stagedQueryState, unstageQueryAction)).toEqual(initialState);
  });

  it('configured staged queries', () => {
    const stagedQueryState = {
      ...initialState,
      stagedQueries: [query],
    };
    const configObject = {
      interval: 3600,
      platform: 'windows',
      logging_type: 'differential',
      query_ids: [query.id],
    };
    const configureStagedQueriesAction = actions.configureStagedQueries(configObject);

    expect(reducer(stagedQueryState, configureStagedQueriesAction)).toEqual({
      stagedQueries: [],
      configuredQueries: [query],
      configurations: [configObject],
    });
  });
});
