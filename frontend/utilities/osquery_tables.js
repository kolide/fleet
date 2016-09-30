import { flatten, map } from 'lodash';
import osqueryTablesJSON from '../osquery_tables';

const appendPlatformKeyToTables = (parsedTables) => {
  return map(parsedTables, (platform) => {
    return platform.tables.map(table => {
      table.platform = platform.key;

      return table;
    });
  });
};

export const normalizeTables = (tablesJSON) => {
  const { tables: parsedTables } = typeof tablesJSON === 'object' ? tablesJSON : JSON.parse(tablesJSON);
  const tablesWithPlatformKey = appendPlatformKeyToTables(parsedTables);

  return flatten(tablesWithPlatformKey);
};

export default normalizeTables(osqueryTablesJSON);
