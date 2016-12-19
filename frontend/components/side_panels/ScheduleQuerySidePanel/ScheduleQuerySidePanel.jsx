import React, { PropTypes } from 'react';

import queryInterface from 'interfaces/query';
import SearchPackQuery from './SearchPackQuery';
import SecondarySidePanelContainer from '../SecondarySidePanelContainer';

const baseClass = 'schedule-query-side-panel';

const ScheduleQuerySidePanel = ({ allQueries, onSelectQuery, selectedQuery }) => {
  return (
    <SecondarySidePanelContainer className={baseClass}>
      <SearchPackQuery
        allQueries={allQueries}
        onSelectQuery={onSelectQuery}
        selectedQuery={selectedQuery}
      />
    </SecondarySidePanelContainer>
  );
};

ScheduleQuerySidePanel.propTypes = {
  allQueries: PropTypes.arrayOf(queryInterface),
  onSelectQuery: PropTypes.func,
  selectedQuery: queryInterface,
};

export default ScheduleQuerySidePanel;
