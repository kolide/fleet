const DEFAULT_NUM_COLUMNS_TO_DISPLAY = 5;

export const columnsToRender = (table, showAllColumns) => {
  if (showAllColumns) return table.columns;

  return table.columns.slice(0, DEFAULT_NUM_COLUMNS_TO_DISPLAY);
};

export const displayTypeForDataType = (dataType) => {
  switch (dataType) {
    case 'TEXT_TYPE':
      return 'text';
    case 'BIGINT_TYPE':
      return 'big int';
    case 'INTEGER_TYPE':
      return 'integer';
    default:
      return dataType;
  }
};

export const shouldShowAllColumns = (table) => {
  const { columns } = table;

  return columns.length <= DEFAULT_NUM_COLUMNS_TO_DISPLAY;
};

export const numAdditionalColumns = (table) => {
  const { columns } = table;

  return columns.length - DEFAULT_NUM_COLUMNS_TO_DISPLAY;
};
