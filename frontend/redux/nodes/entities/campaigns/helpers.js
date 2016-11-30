export const destroyFunc = (campaign) => {
  return Promise.resolve(campaign);
};

export const updateFunc = (campaign, socketData) => {
  return new Promise((resolve, reject) => {
    const { type, data } = socketData;

    if (type === 'totals') {
      return resolve({
        ...campaign,
        totals: data,
      });
    }

    if (type === 'result') {
      const queryResults = campaign.query_results || [];

      return resolve({
        ...campaign,
        query_results: [
          ...queryResults,
          data,
        ],
      });
    }

    return reject();
  });
};

export default { destroyFunc, updateFunc };
