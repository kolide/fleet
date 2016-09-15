import React, { PropTypes } from 'react';
import componentStyles from './styles';

const {
  boxStyles,
  containerStyles,
  headerStyles,
  smallTabStyles,
  tabStyles,
  textStyles,
} = componentStyles;

const StackedWhiteBoxes = ({ headerText, leadText, children }) => {
  return (
    <div style={containerStyles}>
      <div style={smallTabStyles} />
      <div style={tabStyles} />
      <div style={boxStyles}>
        <p style={headerStyles}>{headerText}</p>
        <p style={textStyles}>{leadText}</p>
        {children}
      </div>
    </div>
  );
};

StackedWhiteBoxes.propTypes = {
  headerText: PropTypes.string,
  leadText: PropTypes.string,
  children: PropTypes.element,
};

export default StackedWhiteBoxes;
