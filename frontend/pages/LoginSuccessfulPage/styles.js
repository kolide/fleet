import styles from '../../styles';

const { color, font, padding } = styles;

export default {
  containerStyles: {
    alignItems: 'center',
    display: 'flex',
    flexDirection: 'column',
    paddingTop: '10%',
  },
  loginSuccessStyles: {
    color: color.green,
    textTransform: 'uppercase',
    fontSize: font.large,
    letterSpacing: '2px',
  },
  subtextStyles: {
    fontSize: font.medium,
    color: color.lightGrey,
  },
  whiteBoxStyles: {
    backgroundColor: color.white,
    boxShadow: '0 0 30px 0 rgba(0,0,0,0.30)',
    marginTop: padding.base,
    padding: padding.base,
    paddingTop: padding.most,
    textAlign: 'center',
    width: '384px',
  },
};

