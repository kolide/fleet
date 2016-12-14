import React, { Component, PropTypes } from 'react';
import moment from 'moment';

import Checkbox from 'components/forms/fields/Checkbox';
import Icon from 'components/Icon';
import { platformIconClass } from 'utilities/icon_class';
import queryInterface from 'interfaces/query';

const baseClass = 'queries-list-item';

class QueriesListItem extends Component {
  static propTypes = {
    checked: PropTypes.bool,
    disabled: PropTypes.bool,
    onSelect: PropTypes.func.isRequired,
    query: queryInterface.isRequired,
  };

  renderCheckbox = () => {
    const { checked, disabled, onSelect, query } = this.props;

    return (
      <Checkbox
        checked={checked}
        disabled={disabled}
        name={query.name}
        onChange={onSelect}
      />
    );
  }

  render () {
    const { query } = this.props;
    const { renderCheckbox } = this;
    const updatedTimeAgo = moment(query.updated_at).fromNow();

    return (
      <tr>
        <td>{renderCheckbox()}</td>
        <td>{query.name}</td>
        <td>{query.description}</td>
        <td><Icon name={platformIconClass(query.platform)} /></td>
        <td>{query.author_name}</td>
        <td>{updatedTimeAgo}</td>
      </tr>
    );
  }
}

export default QueriesListItem;

