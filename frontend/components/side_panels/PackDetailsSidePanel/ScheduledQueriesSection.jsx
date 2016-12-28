import React, { Component, PropTypes } from 'react';
import { take } from 'lodash';

import Button from 'components/buttons/Button';
import Icon from 'components/icons/Icon';
import scheduledQueryInterface from 'interfaces/scheduled_query';

const DEFAULT_NUM_QUERIES = 6;

class ScheduledQueriesSection extends Component {
  static propTypes = {
    scheduledQueries: PropTypes.arrayOf(scheduledQueryInterface),
  };

  constructor (props) {
    super(props);

    this.state = { showAllQueries: false };
  }

  onShowAll = () => {
    this.setState({ showAllQueries: true });

    return false;
  }

  renderShowMoreQueries = () => {
    const { showAllQueries } = this.state;
    const scheduledQueryCount = this.props.scheduledQueries.length;
    const shouldRenderShowMore = !showAllQueries && scheduledQueryCount > DEFAULT_NUM_QUERIES;

    if (shouldRenderShowMore) {
      const { onShowAll } = this;
      const numMoreQueries = scheduledQueryCount - DEFAULT_NUM_QUERIES;
      const queryText = numMoreQueries === 1 ? 'Query' : 'Queries';

      return (
        <div>
          <span>{numMoreQueries} More {queryText}</span>
          <Button onClick={onShowAll} variant="unstyled">SHOW</Button>
        </div>
      );
    }

    return false;
  }

  render () {
    const { renderShowMoreQueries } = this;
    const { scheduledQueries } = this.props;
    const { showAllQueries } = this.state;
    const queriesToRender = showAllQueries ? scheduledQueries : take(scheduledQueries, DEFAULT_NUM_QUERIES);

    return (
      <div>
        <p>Queries</p>
        {queriesToRender.map((scheduledQuery) => {
          return (
            <div key={`scheduled-query-${scheduledQuery.id}`}>
              <Icon name="query" /> {scheduledQuery.name}
            </div>
          );
        })}
        {renderShowMoreQueries()}
      </div>
    );
  }
}

export default ScheduledQueriesSection;
