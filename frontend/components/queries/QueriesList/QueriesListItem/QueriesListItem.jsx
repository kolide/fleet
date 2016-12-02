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
    configured: PropTypes.bool,
    onSelect: PropTypes.func.isRequired,
    query: queryInterface.isRequired,
  };

  constructor (props) {
    super(props);

    this.state = { hover: false };
  }

  onHover = (hover) => {
    return () => {
      this.setState({ hover });
    };
  };

  renderCheckbox = () => {
    const { checked, configured, onSelect, query } = this.props;
    const { hover } = this.state;
    const { onHover } = this;

    if (configured) {
      return (
        <span onMouseOver={onHover(true)} onMouseOut={onHover(false)}>
          {!hover && <i className={`${baseClass}__check-icon kolidecon-success-check`} />}
          {hover && <i className={`${baseClass}__check-icon kolidecon-filter`} />}
        </span>
      );
    }

    return (
      <Checkbox
        onChange={onSelect}
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
        <td><Icon name={platformIconClass(query.platform)} /></td>
        <td />
        <td>{updatedTime.fromNow()}</td>
      </tr>
    );
  }
}

export default QueriesListItem;

