import Styles from '../../../../styles';

const { color, padding } = Styles;

export default {
  btnStyle: {
    paddingTop: 0,
    paddingBottom: 0,
    paddingLeft: padding.medium,
    paddingRight: padding.medium,
  },
  hostBtnStyle: {
    backgroundColor: color.brandUltralight,
    boxShadow: `0 3px 0 ${color.brandDark}`,
    color: color.brandDark,
  },
  labelBtnStyle: {
    color: color.accentText,
    backgroundColor: color.bgMedium,
    borderBottom: `1px solid ${color.accentText}`,
    borderLeft: `1px solid ${color.accentText}`,
    borderRight: `1px solid ${color.accentText}`,
    borderTop: `1px solid ${color.accentText}`,
    boxShadow: `0 3px 0 ${color.accentMedium}`,
  },
};
