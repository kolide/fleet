import Styles from '../../../styles';

const { color, font, padding } = Styles;

export default {
  containerStyles: {
    backgroundColor: color.white,
    paddingBottom: padding.base,
    paddingLeft: padding.base,
    paddingRight: padding.base,
    paddingTop: padding.base,
  },
  headerStyles: {
    color: color.textMedium,
  },
  headerHostsTitleStyles: {
    fontSize: font.large,
    marginLeft: padding.half,
    marginRight: padding.base,
    textTransform: 'uppercase',
  },
  headerHostsCountStyles: {
    fontSize: font.small,
  },
};
