import React, { Component, PropTypes } from 'react';
import classNames from 'classnames';
import { isEqual } from 'lodash';
import moment from 'moment';

import Checkbox from 'components/forms/fields/Checkbox';
import Icon from 'components/icons/Icon';
import packInterface from 'interfaces/pack';

const baseClass = 'packs-list-row';

class Row extends Component {
  static propTypes = {
    checked: PropTypes.bool,
    onCheck: PropTypes.func,
    pack: packInterface.isRequired,
  };

  shouldComponentUpdate (nextProps) {
    return !isEqual(this.props, nextProps);
  }

  handleChange = (shouldCheck) => {
    const { onCheck, pack } = this.props;

    return onCheck(shouldCheck, pack.id);
  }

  renderStatusData = () => {
    const { disabled } = this.props.pack;

    const iconClassName = classNames(`${baseClass}__status-icon`, {
      [`${baseClass}__status-icon--enabled`]: !disabled,
      [`${baseClass}__status-icon--disabled`]: disabled,
    });

    if (disabled) {
      return <td><Icon className={iconClassName} name="offline" /> Disabled</td>;
    }

    return <td><Icon className={iconClassName} name="success-check" /> Enabled</td>;
  }

  render () {
    const { checked, pack } = this.props;
    const { handleChange, renderStatusData } = this;
    const updatedTime = moment(pack.updated_at);

    return (
      <tr>
        <td>
          <Checkbox
            name={`select-pack-${pack.id}`}
            onChange={handleChange}
            value={checked}
          />
        </td>
        <td>{pack.name}</td>
        <td>{pack.query_count}</td>
        {renderStatusData()}
        <td />
        <td>{updatedTime.fromNow()}</td>
      </tr>
    );
  }
}

export default Row;

