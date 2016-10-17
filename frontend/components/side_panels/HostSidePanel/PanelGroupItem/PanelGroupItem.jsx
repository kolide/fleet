import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';

import { iconClassForLabel } from './helpers';

const baseClass = 'panel-group-item';

class PanelGroupItem extends Component {
  static propTypes = {
    item: PropTypes.shape({
      hosts_count: PropTypes.number,
      title: PropTypes.string,
      type: PropTypes.string,
    }).isRequired,
    onLabelClick: PropTypes.func,
    isSelected: PropTypes.bool,
  };

  render () {
    const { item, onLabelClick, isSelected } = this.props;
    const {
      hosts_count: count,
      title,
    } = item;
    const wrapperClassName = classnames(baseClass, `${baseClass}__wrapper`, {
      [`${baseClass}__wrapper--is-selected`]: isSelected,
    });

    return (
      <button className={`${wrapperClassName} button button--unstyled`} onClick={onLabelClick}>
        <i className={iconClassForLabel(item)} />
        <span>{title}</span>
        <span>{count}</span>
      </button>
    );
  }
}

export default PanelGroupItem;
