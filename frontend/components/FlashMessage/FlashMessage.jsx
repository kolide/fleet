import React, { PropTypes } from 'react';
import classnames from 'classnames';

import notificationInterface from 'interfaces/notification';
import Icon from 'components/icons/Icon';
import Button from 'components/buttons/Button';

const baseClass = 'flash-message';

const FlashMessage = ({ fullWidth, notification, onRemoveFlash, onUndoActionClick }) => {
  const { alertType, isVisible, message, undoAction } = notification;
  const klass = classnames(baseClass, `${baseClass}--${alertType}`, {
    [`${baseClass}--full-width`]: fullWidth,
  });

  if (!isVisible) {
    return false;
  }

  return (
    <div className={klass}>
      <div className={`${baseClass}__content`}>
        {message}
      </div>
      <div className={`${baseClass}__action`}>
        {undoAction &&
          <Button
            className={`${baseClass}__undo`}
            variant="unstyled"
            onClick={onUndoActionClick(undoAction)}
          >
            Undo
          </Button>
        }

        <Button
          className={`${baseClass}__remove ${baseClass}__remove--${alertType}`}
          variant="unstyled"
          onClick={onRemoveFlash}
        >
          <Icon name="x" />
        </Button>
      </div>
    </div>
  );
};

FlashMessage.propTypes = {
  fullWidth: PropTypes.bool,
  notification: notificationInterface,
  onRemoveFlash: PropTypes.func,
  onUndoActionClick: PropTypes.func,
};

export default FlashMessage;
