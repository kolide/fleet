import React, { Component, PropTypes } from 'react';

import Button from '../../buttons/Button';
import targetInterface from '../../../interfaces/target';
import TargetInfoModal from '../../modals/TargetInfoModal';

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

  handleSelectFromModal = (evt) => {
    const { handleSelect } = this;
    const { onRemoveMoreInfoTarget } = this.props;

    handleSelect(evt);
    onRemoveMoreInfoTarget();
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
        <Button className={`${classBlock}__btn`} text="ADD" onClick={handleSelect} />
        <Button className={`${classBlock}__more-info`} onClick={onMoreInfoClick(target)} text="more info" variant="unstyled" />
      </div>
    );
  }

  renderLabel = () => {
    const { handleSelect, targetIconClass } = this;
    const { onMoreInfoClick, target } = this.props;
    const { count, label } = target;

    return (
      <div className={`${classBlock}__wrapper`}>
        <i className={`${targetIconClass()} ${classBlock}__target-icon`} />
        <span className={`${classBlock}__label-label`}>{label}</span>
        <span className={`${classBlock}__delimeter`}>&bull;</span>
        <span className={`${classBlock}__count`}>{count} hosts</span>
        <Button className={`${classBlock}__btn`} text="ADD" onClick={handleSelect} />
        <Button className={`${classBlock}__more-info`} onClick={onMoreInfoClick(target)} text="more info" variant="unstyled" />
      </div>
    );
  }

  renderTargetInfoModal = () => {
    const { onRemoveMoreInfoTarget, shouldShowModal, target } = this.props;

    if (!shouldShowModal) return false;

    const { handleSelectFromModal } = this;

    return (
      <TargetInfoModal
        className={`${classBlock}__modal-wrapper`}
        onAdd={handleSelectFromModal}
        onExit={onRemoveMoreInfoTarget}
        target={target}
      />
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
