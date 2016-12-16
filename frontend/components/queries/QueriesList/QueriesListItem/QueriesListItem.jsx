import React, { Component, PropTypes } from 'react';
import moment from 'moment';

import Checkbox from 'components/forms/fields/Checkbox';
import Icon from 'components/Icon';
import { isEqual } from 'lodash';
import { platformIconClass } from 'utilities/icon_class';
import scheduledQueryInterface from 'interfaces/scheduled_query';

class QueriesListItem extends Component {
  static propTypes = {
    checked: PropTypes.bool,
    disabled: PropTypes.bool,
    onSelect: PropTypes.func.isRequired,
    scheduledQuery: scheduledQueryInterface.isRequired,
  };

  shouldComponentUpdate (nextProps) {
    if (isEqual(nextProps, this.props)) {
      return false;
    }

    return true;
  }

  onCheck = (value) => {
    const { onSelect, scheduledQuery } = this.props;

    return onSelect(value, scheduledQuery.id);
  }

  loggingTypeString = () => {
    const { scheduledQuery: { snapshot, removed } } = this.props;

    if (snapshot) {
      return 'snapshot';
    }

    if (removed) {
      return 'differential (ignore removes)';
    }

    return 'differential';
  }

  render () {
    const { checked, disabled, scheduledQuery } = this.props;
    const { onCheck } = this;
    const { id, name, interval, platform, updated_at: updatedAt, version } = scheduledQuery;
    const { loggingTypeString } = this;
    const updatedTimeAgo = moment(updatedAt).fromNow();

    return (
      <tr>
        <td>
          <Checkbox
            disabled={disabled}
            name={`scheduled-query-checkbox-${id}`}
            onChange={onCheck}
            value={checked}
          />
        </td>
        <td>{name}</td>
        <td>{interval}</td>
        <td><Icon name={platformIconClass(platform)} /></td>
        <td>{version}</td>
        <td>{loggingTypeString()}</td>
        <td />
        <td>{updatedTimeAgo}</td>
      </tr>
    );
  }
}

export default QueriesListItem;

