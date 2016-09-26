import styles from '../../../styles';

const { border, color, font, padding } = styles;

export default (variant) => {
  const inverse = variant === 'inverse';
  const backgroundImage = inverse ? 'none' : 'linear-gradient(134deg, #7166D9 0%, #C86DD7 100%)';
  const backgroundColor = inverse ? color.white : 'transparent';
  const boxShadow = inverse ? `0 3px 0 ${color.brandLight}` : '0 3px 0 #734893';
  const border = inverse ? '1px solid #EDD6FF' : 'none';

  return {
    backgroundImage,
    backgroundColor,
    borderBottom: border,
    borderLeft: border,
    borderRight: border,
    borderTop: border,
    borderBottomLeftRadius: border.radius.base,
    borderBottomRightRadius: border.radius.base,
    borderTopLeftRadius: border.radius.base,
    borderTopRightRadius: border.radius.base,
    boxShadow,
    boxSizing: 'border-box',
    color: inverse ? color.brand : color.white,
    cursor: 'pointer',
    fontSize: font.large,
    fontWeight: '300',
    letterSpacing: '4px',
    paddingBottom: padding.medium,
    paddingLeft: padding.medium,
    paddingRight: padding.medium,
    paddingTop: padding.medium,
    position: 'relative',
    textTransform: 'uppercase',
    width: '100%',
    ':active': {
      boxShadow: '0 1px 0 #734893, 0 -2px 0 #D1D9E9',
      top: '2px',
    },
    ':focus': {
      outline: 'none',
    },
  };
};
