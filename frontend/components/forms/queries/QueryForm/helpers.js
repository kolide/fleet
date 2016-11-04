export const formNotChanged = (formData, query) => {
  return formData.name === query.name &&
    formData.description === query.description &&
    formData.queryText === query.query;
};

export default { formNotChanged };
