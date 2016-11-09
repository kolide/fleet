import React from 'react';
import classNames from 'classnames';
import { filter, includes, isEqual, noop } from 'lodash';

import { shouldShowModal } from './helpers';
import TargetOption from '../TargetOption';

const SelectTargetsMenuWrapper = (onMoreInfoClick, onRemoveMoreInfoTarget, moreInfoTarget) => {
  const SelectTargetsMenu = ({
    focusedOption,
    instancePrefix,
    onFocus,
    onSelect,
    optionClassName,
    optionComponent,
    options,
    valueArray = [],
    valueKey,
    onOptionRef,
  }) => {
    const Option = optionComponent;
    const renderTargets = (targetType) => {
      const targets = filter(options, { target_type: targetType });

      return targets.map((target, index) => {
        const { disabled: isDisabled } = target;
        const isSelected = includes(valueArray, target);
        const isFocused = isEqual(focusedOption, target);
        const className = classNames(optionClassName, {
          'Select-option': true,
          'is-selected': isSelected,
          'is-focused': isFocused,
          'is-disabled': true,
        });
        const setRef = (ref) => { onOptionRef(ref, isFocused); };
        const isShowModal = shouldShowModal(moreInfoTarget, target);

        return (
          <Option
            className={className}
            instancePrefix={instancePrefix}
            isDisabled={isDisabled}
            isFocused={isFocused}
            isSelected={isSelected}
            key={`option-${index}-${target[valueKey]}`}
            onFocus={onFocus}
            onSelect={noop}
            option={target}
            optionIndex={index}
            ref={setRef}
          >
            <TargetOption
              target={moreInfoTarget && isShowModal ? moreInfoTarget : target}
              onSelect={onSelect}
              onRemoveMoreInfoTarget={onRemoveMoreInfoTarget}
              onMoreInfoClick={onMoreInfoClick}
              shouldShowModal={isShowModal}
            />
          </Option>
        );
      });
    };

    return (
      <div>
        <div>hosts</div>
        {renderTargets('hosts')}
        <div>labels</div>
        {renderTargets('labels')}
      </div>
    );
  };

  return SelectTargetsMenu;
};

export default SelectTargetsMenuWrapper;
