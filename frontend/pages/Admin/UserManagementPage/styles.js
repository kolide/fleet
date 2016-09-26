import Styles from '../../../styles';

const { border, color, font, padding } = Styles;

export default {
  addUserButtonStyles: {
    backgroundColor: color.brand,
    backgroundImage: 'none',
    boxShadow: `0 4px 0 ${color.brandLight}`,
    fontSize: font.medium,
    height: '38px',
    letterSpacing: 'normal',
    marginTop: 0,
    marginLeft: padding.half,
    padding: 0,
    width: '145px',
  },
  addUserWrapperStyles: {
    float: 'right',
  },
  avatarStyles: {
    display: 'block',
    marginLeft: 'auto',
    marginRight: 'auto',
  },
  containerStyles: {
    backgroundColor: color.white,
    minHeight: '100px',
    paddingBottom: '190px',
    paddingRight: padding.most,
    paddingTop: padding.base,
  },
  nameStyles: {
    fontWeight: font.weight.bold,
    lineHeight: '51px',
    margin: 0,
    padding: 0,
  },
  numUsersStyles: {
    borderBottom: '1px solid #EFF0F4',
    color: color.textMedium,
    display: 'inline-block',
    fontSize: font.large,
    marginLeft: padding.most,
    paddingBottom: padding.half,
    width: '260px',
  },
  userHeaderStyles: {
    backgroundColor: color.brand,
    color: color.white,
    height: '51px',
    marginBottom: padding.half,
    textAlign: 'center',
    width: '100%',
  },
  userDetailsStyles: {
    paddingLeft: padding.half,
    paddingRight: padding.half,
  },
  userEmailStyles: {
    fontSize: font.mini,
    color: color.link,
  },
  userLabelStyles: {
    float: 'left',
    fontSize: font.small,
  },
  usernameStyles: {
    color: color.brand,
    fontSize: font.medium,
    textTransform: 'uppercase',
  },
  userPositionStyles: {
    fontSize: font.small,
  },
  userStatusStyles: (enabled) => {
    return {
      color: enabled ? color.success : color.textMedium,
      float: 'right',
      fontSize: font.small,
    };
  },
  userStatusWrapperStyles: {
    borderBottomColor: color.borderMedium,
    borderBottomStyle: 'solid',
    borderBottomWidth: '1px',
    borderTopColor: color.borderMedium,
    borderTopStyle: 'solid',
    borderTopWidth: '1px',
    marginTop: padding.half,
    paddingTop: padding.half,
    paddingBottom: padding.half,
  },
  userWrapperStyles: {
    boxShadow: border.shadow.blur,
    display: 'inline-block',
    height: '390px',
    marginLeft: padding.most,
    marginTop: padding.most,
    width: '239px',
  },
  usersWrapperStyles: {
  },
};
