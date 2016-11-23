import React, { Component, PropTypes } from 'react';
import moment from 'moment';

import Checkbox from 'components/forms/fields/Checkbox';
import { platformIconClass } from 'utilities/icon_class';
import queryInterface from 'interfaces/query';

class QueriesListItem extends Component {
  static propTypes = {
    checked: PropTypes.bool,
    configured: PropTypes.bool,
    onSelect: PropTypes.func.isRequired,
    query: queryInterface.isRequired,
  };

  renderCheckbox = () => {
    const { checked, configured, onSelect, query } = this.props;

    if (configured) {
      return <i className="kolidecon-success-check" />;
    }

    return (
      <Checkbox
        onClick={onSelect}
        name={query.name}
        checked={checked}
      />
    );
  }

  render () {
    const { query } = this.props;
    const { renderCheckbox } = this;
    const updatedTime = moment(query.updated_at);

    return (
      <tr>
        <td>
          {renderCheckbox()}
        </td>
        <td>{query.name}</td>
        <td>{query.description}</td>
        <td><i className={platformIconClass(query.platform)} /></td>
        <td />
        <td>{updatedTime.fromNow()}</td>
      </tr>
    );
  }
}

export default QueriesListItem;

