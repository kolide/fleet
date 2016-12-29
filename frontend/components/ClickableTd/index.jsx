import React, { PropTypes } from 'react';

const ClickableTd = ({ children, className, onClick }) => {
  return <td className={className}><a onClick={onClick} tabIndex={-1}>{children}</a></td>;
};

ClickableTd.propTypes = {
  children: PropTypes.node,
  className: PropTypes.string,
  onClick: PropTypes.func.isRequired,
};

export default ClickableTd;
