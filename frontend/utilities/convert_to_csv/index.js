import { keys } from 'lodash';

const convertToCSV = (objArray) => {
  const fields = keys(objArray[0]);
  const jsonFields = fields.map(field => JSON.stringify(field));
  const rows = objArray.map((row) => {
    return fields.map(field => JSON.stringify(row[field])).join(',');
  });


  rows.unshift(jsonFields.join(','));

  return rows.join('\n');
};

export default convertToCSV;
