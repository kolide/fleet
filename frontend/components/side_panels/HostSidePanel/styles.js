import Styles from '../../../styles';

const { color, padding } = Styles;

export default {
  containerStyles: {
    color: color.textMedium,
    paddingLeft: 0,
    paddingRight: 0,
  },
  hrStyles: {
    color: color.textMedium,
    marginBottom: padding.base,
    marginLeft: padding.auto,
    marginRight: padding.auto,
    marginTop: padding.base,
    width: '80%',
  },
  PanelGroupItemStyles: {
    containerStyles: (selected) => {
      return {
        backgroundColor: selected ? color.brand : color.white,
        color: selected ? color.white : color.textMedium,
        cursor: 'pointer',
        paddingLeft: padding.large,
        paddingRight: padding.large,
        paddingTop: padding.half,
        paddingBottom: padding.half,
      };
    },
    itemStyles: {
      display: 'inline-block',
    },
  },
};
