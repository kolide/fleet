import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import Avatar from '../../../Avatar';
import Button from '../../../buttons/Button';
import componentStyles from '../../../../pages/Admin/UserManagementPage/UserBlock/styles';
import InputFieldWithIcon from '../../fields/InputFieldWithIcon';

class EditUserForm extends Component {
  static propTypes = {
    onSubmit: PropTypes.func,
    user: PropTypes.object,
  };

  constructor (props) {
    super(props);

    this.state = {
      formData: {},
    };
  }

  onInputChange = (fieldName) => {
    return (evt) => {
      this.setState({
        [fieldName]: evt.target.value,
      });

      return false;
    };
  }

  onFormSubmit = (evt) => {
    evt.preventDefault();
    const { formData } = this.state;
    const { onSubmit } = this.props;

    return onSubmit(formData);
  }

  render () {
    const {
      avatarStyles,
      formButtonStyles,
      nameStyles,
      userDetailsStyles,
      userEmailStyles,
      userHeaderStyles,
      userLabelStyles,
      usernameStyles,
      userPositionStyles,
      userStatusStyles,
      userStatusWrapperStyles,
      userWrapperStyles,
    } = componentStyles;
    const { user } = this.props;
    const {
      admin,
      email,
      enabled,
      name,
      position,
      username,
    } = user;
    const { onFormSubmit, onInputChange } = this;

    return (
      <form style={userWrapperStyles} onSubmit={onFormSubmit}>
        <InputFieldWithIcon
          defaultValue={name}
          name="name"
          onChange={onInputChange('name')}
          style={{ marginTop: 0 }}
        />
        <div style={userDetailsStyles}>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Button
              onClick={this.props.onCancel}
              style={formButtonStyles}
              text="Cancel"
              variant="inverse"
            />
            <Button
              style={formButtonStyles}
              text="Submit"
              type="submit"
            />
          </div>
        </div>
      </form>
    );
  }
}

export default radium(EditUserForm);
