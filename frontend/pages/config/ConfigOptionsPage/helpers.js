const configOptionDropdownOptions = (configOptions) => {
  return configOptions.map((option) => {
    return {
      disabled: option.read_only,
      label: option.name,
      value: option.name,
    };
  });
};

export default { configOptionDropdownOptions };
