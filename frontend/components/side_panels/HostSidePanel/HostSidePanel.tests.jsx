import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import HostSidePanel from 'components/side_panels/HostSidePanel';
import { labelStub, statusLabelStub } from 'test/stubs';

describe('HostSidePanel - component', () => {
  const defaultProps = {
    labels: [labelStub],
    onAddLabelClick: noop,
    onAddHostClick: noop,
    onLabelClick: noop,
    statusLabels: { ...statusLabelStub, loading_counts: false, online_count: 10 },
  };

  describe('#shouldComponentUpdate', () => {
    it('does not re-render when only the status label loading changes', () => {
      const SidePanel = mount(<HostSidePanel {...defaultProps} />);
      const sidePanelNode = SidePanel.node;
      const updatedStatusLabels = { ...statusLabelStub, loading_counts: true };
      const updateProps = { ...defaultProps, statusLabels: updatedStatusLabels };

      expect(sidePanelNode.shouldComponentUpdate(updateProps)).toEqual(false);
    });

    it('re-renders when the status label counts change', () => {
      const SidePanel = mount(<HostSidePanel {...defaultProps} />);
      const sidePanelNode = SidePanel.node;
      const updatedStatusLabels = { ...statusLabelStub, loading_counts: true, online_count: 11 };
      const updateProps = { ...defaultProps, statusLabels: updatedStatusLabels };

      expect(sidePanelNode.shouldComponentUpdate(updateProps)).toEqual(true);
    });

    it('re-renders when the status label counts change', () => {
      const SidePanel = mount(<HostSidePanel {...defaultProps} />);
      const sidePanelNode = SidePanel.node;
      const updatedStatusLabels = { ...statusLabelStub, loading_counts: true, online_count: 11 };
      const updateProps = { ...defaultProps, statusLabels: updatedStatusLabels };

      expect(sidePanelNode.shouldComponentUpdate(updateProps)).toEqual(true);
    });

    it('re-renders when the state labelFilter changes', () => {
      const SidePanel = mount(<HostSidePanel {...defaultProps} />);
      const sidePanelNode = SidePanel.node;
      const updatedState = { labelFilter: 'my label' };

      expect(sidePanelNode.shouldComponentUpdate(defaultProps, updatedState)).toEqual(true);
    });
  });
});
