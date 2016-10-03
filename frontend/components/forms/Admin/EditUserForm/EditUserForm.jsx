import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import Avatar from '../../../Avatar';
import Button from '../../../buttons/Button';
import componentStyles from '../../../../pages/Admin/UserManagementPage/UserBlock/styles';
import InputFieldWithLabel from '../../fields/InputFieldWithLabel';

class EditUserForm extends Component {
  static propTypes = {
    onCancel: PropTypes.func,
    onSubmit: PropTypes.func,
    user: PropTypes.object,
  };

  constructor (props) {
    super(props);

    const { user } = props;

    this.state = {
      formData: {
        ...user,
      },
    };
  }

  onInputChange = (fieldName) => {
    return (evt) => {
      const { formData } = this.state;

      this.setState({
        formData: {
          ...formData,
          [fieldName]: evt.target.value,
        },
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
      userWrapperStyles,
    } = componentStyles;
    const { user } = this.props;
    const {
      email,
      name,
      position,
      username,
    } = user;
    const { onFormSubmit, onInputChange } = this;

    return (
      <form style={[userWrapperStyles, { boxSizing: 'border-box', padding: '10px' }]} onSubmit={onFormSubmit}>
        <InputFieldWithLabel
          defaultValue={name}
          label="name"
          name="name"
          onChange={onInputChange('name')}
          style={{ container: { marginTop: 0 } }}
        />
        <Avatar user={user} style={avatarStyles} />
        <InputFieldWithLabel
          defaultValue={username}
          label="username"
          name="username"
          onChange={onInputChange('username')}
          style={{ container: { marginTop: 0 }, input: { color: '#AE6DDF' } }}
        />
        <InputFieldWithLabel
          defaultValue={position}
          label="position"
          name="position"
          onChange={onInputChange('position')}
          style={{ container: { marginTop: 0 } }}
        />
        <InputFieldWithLabel
          defaultValue={email}
          label="email"
          name="email"
          onChange={onInputChange('email')}
          style={{ container: { marginTop: 0 }, input: { color: '#4A90E2' } }}
        />
        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '10px' }}>
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
      </form>
    );
  }
}

export default radium(EditUserForm);
