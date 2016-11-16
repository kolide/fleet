import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';

class Rocker extends Component {

  static propTypes = {
    className: PropTypes.string,
    name: PropTypes.string,
    value: PropTypes.string,
    options: PropTypes.shape({
      aText: PropTypes.string,
      aIcon: PropTypes.string,
      bText: PropTypes.string,
      bIcon: PropTypes.string,
    }),
  };

  render () {
    const { className, name, value } = this.props;
    const baseClass = 'kolide-rocker';

    const rockerClasses = classnames(baseClass, className);

    return (
      <div className={rockerClasses}>
        <label className={`${baseClass}__label`}>
          <input className={`${baseClass}__checkbox`} type="checkbox" value={value} name={name} />
          <span className={`${baseClass}__switch ${baseClass}__switch--opt-b`}>
            <span className={`${baseClass}__text`}><i className="kolidecon kolidecon-th-large"></i> Grid</span>
          </span>
          <span className={`${baseClass}__switch ${baseClass}__switch--opt-a`}>
            <span className={`${baseClass}__text`}><i className="kolidecon kolidecon-th-list"></i> List</span>
          </span>
        </label>
      </div>
    );
  }
}

export default Rocker;
