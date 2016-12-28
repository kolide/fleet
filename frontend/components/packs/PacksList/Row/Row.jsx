import React, { Component, PropTypes } from 'react';
import classNames from 'classnames';
import { isEqual } from 'lodash';
import moment from 'moment';

import Checkbox from 'components/forms/fields/Checkbox';
import ClickableTd from 'components/ClickableTd';
import Icon from 'components/icons/Icon';
import packInterface from 'interfaces/pack';

const baseClass = 'packs-list-row';

class Row extends Component {
  static propTypes = {
    checked: PropTypes.bool,
    onCheck: PropTypes.func,
    onSelect: PropTypes.func,
    pack: packInterface.isRequired,
  };

  shouldComponentUpdate (nextProps) {
    return !isEqual(this.props, nextProps);
  }

  handleChange = (shouldCheck) => {
    const { onCheck, pack } = this.props;

    return onCheck(shouldCheck, pack.id);
  }

  handleSelect = () => {
    const { onSelect, pack } = this.props;

    return onSelect(pack);
  }

  renderStatusData = () => {
    const { disabled } = this.props.pack;
    const { handleSelect } = this;
    const iconClassName = classNames(`${baseClass}__status-icon`, {
      [`${baseClass}__status-icon--enabled`]: !disabled,
      [`${baseClass}__status-icon--disabled`]: disabled,
    });

    if (disabled) {
      return (
        <ClickableTd onClick={handleSelect} tabIndex={-1}>
          <Icon className={iconClassName} name="offline" /> Disabled
        </ClickableTd>
      );
    }

    return (
      <ClickableTd onClick={handleSelect}>
        <Icon className={iconClassName} name="success-check" /> Enabled
      </ClickableTd>
    );
  }

  render () {
    const { checked, pack } = this.props;
    const { handleChange, handleSelect, renderStatusData } = this;
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
        <ClickableTd onClick={handleSelect}>{pack.name}</ClickableTd>
        <ClickableTd onClick={handleSelect}>{pack.query_count}</ClickableTd>
        {renderStatusData()}
        <td />
        <ClickableTd onClick={handleSelect}>{updatedTime.fromNow()}</ClickableTd>
      </tr>
    );
  }
}

export default Row;

