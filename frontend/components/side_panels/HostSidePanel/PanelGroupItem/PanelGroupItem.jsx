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
      count,
      display_text: displayText,
    } = item;
    const wrapperClassName = classnames(
      baseClass,
      'button',
      'button--unstyled',
      `${baseClass}__${item.type.toLowerCase()}`,
      `${baseClass}__${item.type.toLowerCase()}--${displayText.toLowerCase()}`,
      {
        [`${baseClass}--selected`]: isSelected,
      }
    );

    return (
      <button className={wrapperClassName} onClick={onLabelClick}>
        <i className={`${iconClassForLabel(item)} ${baseClass}__icon`} />
        <span className={`${baseClass}__name`}>{displayText}</span>
        <span className={`${baseClass}__count`}>{count}</span>
      </button>
    );
  }
}

export default PanelGroupItem;
