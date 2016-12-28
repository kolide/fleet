import React, { PropTypes } from 'react';

const ClickableTd = ({ children, onClick }) => {
  return <td><a onClick={onClick} tabIndex={-1}>{children}</a></td>;
};

ClickableTd.propTypes = {
  children: PropTypes.node,
  onClick: PropTypes.func.isRequired,
};

export default ClickableTd;
