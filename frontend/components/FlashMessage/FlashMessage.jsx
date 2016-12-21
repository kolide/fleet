import React, { PropTypes } from 'react';

import notificationInterface from 'interfaces/notification';
import Icon from 'components/Icon';
import Button from 'components/buttons/Button';

const baseClass = 'flash-message';

const FlashMessage = ({ notification, onRemoveFlash, onUndoActionClick }) => {
  const { alertType, isVisible, message, undoAction } = notification;

  if (!isVisible) {
    return false;
  }

  const alertIcon = alertType === 'success' ? 'success-check' : 'warning-filled';

  return (
    <div className={`${baseClass} ${baseClass}--${alertType}`}>
      <div className={`${baseClass}__content`}>
        <Icon name={alertIcon} /> <span>{message}</span>

        {undoAction &&
          <Button
            className={`${baseClass}__undo`}
            variant="unstyled"
            onClick={onUndoActionClick(undoAction)}
            text="undo"
          />
        }
      </div>
      <div className={`${baseClass}__action`}>
        <Button
          className={`${baseClass}__remove ${baseClass}__remove--${alertType}`}
          variant="unstyled"
          onClick={onRemoveFlash}
          text={<Icon name="x" />}
        />
      </div>
    </div>
  );
};

FlashMessage.propTypes = {
  notification: notificationInterface,
  onRemoveFlash: PropTypes.func,
  onUndoActionClick: PropTypes.func,
};

export default FlashMessage;
