import React, { Component, PropTypes } from 'react';
import radium from 'radium';

import { calculateTooltipDirection } from './helpers';
import ClickOutside from '../../ClickOutside';
import componentStyles from './styles';

export class ElipsisMenu extends Component {
  static propTypes = {
    children: PropTypes.node,
    positionStyles: PropTypes.object,
  };

  constructor (props) {
    super(props);

    this.state = {
      showChildren: false,
    };
  }

  componentDidMount () {
    const { setTooltipDirection } = this;

    global.window.addEventListener('resize', setTooltipDirection);

    return setTooltipDirection();
  }

  componentWillUnmount () {
    const { setTooltipDirection } = this;

    global.window.removeEventListener('resize', setTooltipDirection);

    return false;
  }

  onToggleChildren = () => {
    const { showChildren } = this.state;

    this.setState({ showChildren: !showChildren });

    return false;
  }

  setElement = (element) => {
    this.element = element;
  }

  setTooltipDirection = () => {
    if (this.element) {
      const tooltipDirection = calculateTooltipDirection(this.element);

      this.setState({ tooltipDirection });
    }

    return false;
  }

  renderChildren = () => {
    const { children } = this.props;
    const { childrenWrapperStyles } = componentStyles;
    const { showChildren, tooltipDirection } = this.state;
    const triangleDirection = tooltipDirection === 'left' ? 'right' : 'left';

    if (!showChildren) return false;

    return (
      <div
        className={`container-triangle ${triangleDirection}`}
        style={childrenWrapperStyles(tooltipDirection)}
      >
        {children}
      </div>
    );
  }

  render () {
    const { containerStyles, elipsisStyles } = componentStyles;
    const { onToggleChildren, renderChildren, setElement } = this;
    const { positionStyles } = this.props;

    return (
      <div
        onClick={onToggleChildren}
        ref={setElement}
        style={[containerStyles, positionStyles]}
      >
        <span style={elipsisStyles}>&bull; &bull; &bull;</span>
        {renderChildren()}
      </div>
    );
  }
}

const StyledComponent = radium(ElipsisMenu);
export default ClickOutside(StyledComponent, {
  onOutsideClick: (component) => {
    return (evt) => {
      evt.preventDefault();
      component.setState({ showChildren: false });

      return false;
    };
  },
});
