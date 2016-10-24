import React from 'react';

import Button from '../../../buttons/Button';
import componentStyles from './styles';
import targetInterface from '../../../../interfaces/target';

const classBlock = 'target-option';

const TargetOption = ({ count, ip, label, platform, target_type: targetType }) => {
  const { btnStyle, hostBtnStyle, labelBtnStyle } = componentStyles;
  const iconClass = platform === 'darwin' ? 'kolidecon-apple' : `kolidecon-${platform}`;

  if (targetType === 'hosts') {
    return (
      <div className={`${classBlock}-wrapper`}>
        <Button style={[btnStyle, hostBtnStyle]} text="HOST" />
        <i className={`${classBlock}__icon ${iconClass}`} />
        <span className={`${classBlock}__label-host`}>{label}</span>
        &bull;
        <span className={`${classBlock}__ip`}>{ip}</span>
      </div>
    );
  }

  return (
    <div className={`${classBlock}-wrapper`}>
      <Button style={[btnStyle, labelBtnStyle]} text="LABEL" />
      <span className={`${classBlock}__label-label`}>{label}</span>
      &bull;
      <span className={`${classBlock}__count`}>{count} hosts</span>
    </div>
  );
};

TargetOption.propTypes = targetInterface;

export default TargetOption;
