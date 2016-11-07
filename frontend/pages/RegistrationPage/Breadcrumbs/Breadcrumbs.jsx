import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';

class Breadcrumbs extends Component {
  static propTypes = {
    onClick: PropTypes.func,
    page: PropTypes.number,
  };

  static defaultProps = {
    page: 1,
  };

  onClick = (page) => {
    return (evt) => {
      evt.preventDefault();

      const { onClick: handleClick } = this.props;

      return handleClick(page);
    };
  }

  render () {
    const { onClick } = this;
    const { page } = this.props;
    const page1ClassName = classnames('button--unstyled', 'page-1-btn', {
      'is-active': page >= 1,
    });
    const page2ClassName = classnames('button--unstyled', 'page-2-btn', {
      'is-active': page >= 2,
    });
    const page3ClassName = classnames('button--unstyled', 'page-3-btn', {
      'is-active': page >= 3,
    });

    return (
      <div>
        <button className={page1ClassName} onClick={onClick(1)} />
        <button className={page2ClassName} onClick={onClick(2)} />
        <button className={page3ClassName} onClick={onClick(3)} />
      </div>
    );
  }
}

export default Breadcrumbs;
