import React, { Component, PropTypes } from 'react';

import Button from '../../../buttons/Button';
import componentStyles from './styles';
import Modal from '../../../../components/Modal';
import targetInterface from '../../../../interfaces/target';

const classBlock = 'target-option';

class TargetOption extends Component {
  static propTypes = {
    onMoreInfoClick: PropTypes.func,
    onRemoveMoreInfoTarget: PropTypes.func,
    onSelect: PropTypes.func,
    shouldShowModal: PropTypes.bool,
    target: targetInterface.isRequired,
  };

  handleSelect = (evt) => {
    const { onSelect, target } = this.props;
    return onSelect(target, evt);
  }

  hostPlatformIconClass = () => {
    const { platform } = this.props.target;

    return platform === 'darwin' ? 'kolidecon-apple' : `kolidecon-${platform}`;
  }

  targetIconClass = () => {
    const { label, target_type: targetType } = this.props.target;

    if (label.toLowerCase() === 'all hosts') {
      return 'kolidecon-all-hosts';
    }

    if (targetType === 'hosts') {
      return 'kolidecon-single-host';
    }

    return 'kolidecon-label';
  }

  renderHost = () => {
    const { btnStyle } = componentStyles;
    const { handleSelect, hostPlatformIconClass, targetIconClass } = this;
    const { onMoreInfoClick, target } = this.props;
    const { ip, label } = target;

    return (
      <div className={`${classBlock}__wrapper`}>
        <i className={`${targetIconClass()} ${classBlock}__target-icon`} />
        <i className={`${classBlock}__icon ${hostPlatformIconClass()}`} />
        <span className={`${classBlock}__label-host`}>{label}</span>
        <span className={`${classBlock}__delimeter`}>&bull;</span>
        <span className={`${classBlock}__ip`}>{ip}</span>
        <Button style={btnStyle} text="ADD" onClick={handleSelect} />
        <button className={`btn--unstyled ${classBlock}__more-info`} onClick={onMoreInfoClick(target)}>more info</button>
      </div>
    );
  }

  renderLabel = () => {
    const { btnStyle } = componentStyles;
    const { handleSelect, targetIconClass } = this;
    const { onMoreInfoClick, target } = this.props;
    const { count, label } = target;

    return (
      <div className={`${classBlock}__wrapper`}>
        <i className={`${targetIconClass()} ${classBlock}__target-icon`} />
        <span className={`${classBlock}__label-label`}>{label}</span>
        <span className={`${classBlock}__delimeter`}>&bull;</span>
        <span className={`${classBlock}__count`}>{count} hosts</span>
        <Button style={btnStyle} text="ADD" onClick={handleSelect} />
        <button className={`btn--unstyled ${classBlock}__more-info`} onClick={onMoreInfoClick(target)}>more info</button>
      </div>
    );
  }

  renderTargetInfoModal = () => {
    const { onRemoveMoreInfoTarget, shouldShowModal, target } = this.props;

    const { label } = target;

    if (!shouldShowModal) return false;

    return (
      <Modal
        className={`${classBlock}__target-modal`}
        onExit={onRemoveMoreInfoTarget}
        title={label}
      >
        <p>FIXME</p>
      </Modal>
    );
  }

  render () {
    const { target_type: targetType } = this.props.target;
    const { renderHost, renderLabel, renderTargetInfoModal } = this;

    if (targetType === 'hosts') {
      return (
        <div>
          {renderHost()}
          {renderTargetInfoModal()}
        </div>
      );
    }

    return (
      <div>
        {renderLabel()}
        {renderTargetInfoModal()}
      </div>
    );
  }
}

export default TargetOption;
