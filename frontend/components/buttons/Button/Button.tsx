import * as React from 'react';
import classnames from 'classnames';

const baseClass = 'button';

interface IButtonProps {
  className: string;
  disabled: boolean;
  onClick: (evt: any) => boolean;
  text: string;
  type: string;
  variant: string;
}

class Button extends React.Component<IButtonProps, any> {
  static defaultProps = {
    variant: 'default',
  };

  handleClick = (evt: any) => {
    const { disabled, onClick } = this.props;

    if (disabled) {
      return false;
    }

    if (onClick) {
      onClick(evt);
    }

    return false;
  }

  render () {
    const { handleClick } = this;
    const { className, disabled, text, type, variant } = this.props;
    const fullClassName = classnames(`${baseClass}__${variant}`, className, {
      [baseClass]: variant !== 'unstyled',
      [`${baseClass}__${variant}--disabled`]: disabled,
    });

    return (
      <button
        className={fullClassName}
        disabled={disabled}
        onClick={handleClick}
        type={type}
      >
        {text}
      </button>
    );
  }
}

export default Button;
