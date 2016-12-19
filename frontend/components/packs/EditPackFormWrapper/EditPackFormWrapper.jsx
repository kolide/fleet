import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import EditPackForm from 'components/forms/packs/EditPackForm';
import packInterface from 'interfaces/pack';

class EditPackFormWrapper extends Component {
  static propTypes = {
    className: PropTypes.string,
    handleSubmit: PropTypes.func,
    pack: packInterface.isRequired,
  };

  constructor (props) {
    super(props);

    this.state = { isEditing: false };
  }

  onEditPack = (evt) => {
    evt.preventDefault();

    this.setState({ isEditing: true });

    return false;
  }

  render () {
    const { isEditing } = this.state;
    const { onEditPack } = this;
    const { className, handleSubmit, pack } = this.props;

    if (isEditing) {
      return (
        <EditPackForm
          className={className}
          formData={pack}
          handleSubmit={handleSubmit}
        />
      );
    }

    return (
      <div className={className}>
        <Button
          onClick={onEditPack}
          text="EDIT"
          type="button"
          variant="brand"
        />
      </div>
    );
  }
}

export default EditPackFormWrapper;
