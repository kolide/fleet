import expect from 'expect';

import helpers from './helpers';

const campaign = {
  created_at: '0001-01-01T00:00:00Z',
  deleted: false,
  deleted_at: null,
  id: 4,
  query_id: 12,
  status: 0,
  updated_at: '0001-01-01T00:00:00Z',
  user_id: 1,
};
const campaignWithResults = {
  created_at: '0001-01-01T00:00:00Z',
  deleted: false,
  deleted_at: null,
  id: 4,
  query_id: 12,
  query_results: [
    { distributed_query_execution_id: 4, hosts: [], rows: [] },
  ],
  status: 0,
  totals: {
    count: 3,
    online: 2,
  },
  updated_at: '0001-01-01T00:00:00Z',
  user_id: 1,
};
const { destroyFunc, updateFunc } = helpers;
const resultSocketData = {
  type: 'result',
  data: {
    distributed_query_execution_id: 5,
    hosts: [],
    rows: [],
  },
};
const totalsSocketData = {
  type: 'totals',
  data: {
    count: 5,
    online: 1,
  },
};

describe('campaign entity - helpers', () => {
  describe('#destroyFunc', () => {
    it('returns the campaign', (done) => {
      destroyFunc(campaign)
        .then((response) => {
          expect(response).toEqual(campaign);
          done();
        })
        .catch(done);
    });
  });

  describe('#updateFunc', () => {
    it('appends query results to the campaign when the campaign has query results', (done) => {
      updateFunc(campaignWithResults, resultSocketData)
        .then((response) => {
          expect(response.query_results).toEqual([
            { distributed_query_execution_id: 4, hosts: [], rows: [] },
            resultSocketData.data,
          ]);
          done();
        })
        .catch(done);
    });

    it('adds query results to the campaign when the campaign does not have query results', (done) => {
      updateFunc(campaign, resultSocketData)
        .then((response) => {
          expect(response.query_results).toEqual([
            resultSocketData.data,
          ]);
          done();
        })
        .catch(done);
    });

    it('updates totals on the campaign when the campaign has totals', (done) => {
      updateFunc(campaignWithResults, totalsSocketData)
        .then((response) => {
          expect(response.totals).toEqual(totalsSocketData.data);
          done();
        })
        .catch(done);
    });

    it('adds totals to the campaign when the campaign does not have totals', (done) => {
      updateFunc(campaign, totalsSocketData)
        .then((response) => {
          expect(response.totals).toEqual(totalsSocketData.data);
          done();
        })
        .catch(done);
    });
  });
});
