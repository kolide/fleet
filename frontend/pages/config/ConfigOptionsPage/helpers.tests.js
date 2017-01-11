import expect from 'expect';

import { configOptionStub } from 'test/stubs';
import helpers from 'pages/config/ConfigOptionsPage/helpers';

describe('ConfigOptionsPage - helpers', () => {
  describe('#configOptionDropdownOptions', () => {
    const configOptions = [
      configOptionStub,
      { ...configOptionStub, id: 2, name: 'another_config_option' },
      { ...configOptionStub, id: 3, name: 'third_config_option', read_only: true },
    ];

    it('returns the available dropdown options', () => {
      expect(helpers.configOptionDropdownOptions(configOptions)).toEqual([
        { label: configOptionStub.name, value: configOptionStub.name, disabled: false },
        { label: 'another_config_option', value: 'another_config_option', disabled: false },
        { label: 'third_config_option', value: 'third_config_option', disabled: true },
      ]);
    });
  });
});
