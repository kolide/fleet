import Styles from '../../../styles';

const { color, font, padding } = Styles;

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
  containerStyles: {
    backgroundColor: color.white,
    minHeight: '100px',
    paddingBottom: '190px',
    paddingRight: padding.most,
    paddingTop: padding.base,
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
  usersWrapperStyles: {
  },
};
