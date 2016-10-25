import Styles from '../../../styles';

const { color, padding } = Styles;

export default {
  btnStyle: {
    backgroundColor: color.brand,
    borderBottom: 'none',
    borderLeft: 'none',
    borderRight: 'none',
    borderTop: 'none',
    boxShadow: `0 3px 0 ${color.brandDark}`,
    color: color.white,
    float: 'right',
    paddingTop: 0,
    paddingBottom: 0,
    paddingLeft: padding.medium,
    paddingRight: padding.medium,
  },
};
