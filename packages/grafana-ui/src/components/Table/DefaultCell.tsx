import { cx } from '@emotion/css';
import React, { ReactElement, useState } from 'react';

import { DisplayValue, formattedValueToString } from '@grafana/data';
import { TableCellDisplayMode } from '@grafana/schema';

import { useStyles2 } from '../../themes';
import { getCellLinks } from '../../utils';
import { clearLinkButtonStyles } from '../Button';
import { DataLinksContextMenu } from '../DataLinks/DataLinksContextMenu';

import { CellActions } from './CellActions';
import { TableStyles } from './styles';
import { TableCellProps, CustomCellRendererProps, TableCellOptions } from './types';
import { getCellColors, getCellOptions } from './utils';

export const DefaultCell = (props: TableCellProps) => {
  const { field, cell, tableStyles, row, cellProps, frame } = props;

  const inspectEnabled = Boolean(field.config.custom?.inspect);
  const displayValue = field.display!(cell.value);

  const showFilters = props.onCellFilterAdded && field.config.filterable;
  const showActions = (showFilters && cell.value !== undefined) || inspectEnabled;
  const cellOptions = getCellOptions(field);
  const hasLinks = Boolean(getCellLinks(field, row)?.length);
  const clearButtonStyle = useStyles2(clearLinkButtonStyles);
  const [hover, setHover] = useState(false);
  let value: string | ReactElement;

  const onMouseLeave = () => {
    setHover(false);
  };
  const onMouseEnter = () => {
    setHover(true);
  };

  if (cellOptions.type === TableCellDisplayMode.Custom) {
    const CustomCellComponent: React.ComponentType<CustomCellRendererProps> = cellOptions.cellComponent;
    value = <CustomCellComponent field={field} value={cell.value} rowIndex={row.index} frame={frame} />;
  } else {
    if (React.isValidElement(cell.value)) {
      value = cell.value;
    } else {
      value = formattedValueToString(displayValue);
    }
  }

  const isStringValue = typeof value === 'string';

  const cellStyle = getCellStyle(tableStyles, cellOptions, displayValue, inspectEnabled, isStringValue);

  if (isStringValue) {
    let justifyContent = cellProps.style?.justifyContent;

    if (justifyContent === 'flex-end') {
      cellProps.style = { ...cellProps.style, textAlign: 'right' };
    } else if (justifyContent === 'center') {
      cellProps.style = { ...cellProps.style, textAlign: 'center' };
    }
  }

  return (
    <div
      {...cellProps}
      onMouseEnter={showActions ? onMouseEnter : undefined}
      onMouseLeave={showActions ? onMouseLeave : undefined}
      className={cellStyle}
    >
      {!hasLinks && (isStringValue ? `${value}` : <div className={tableStyles.cellText}>{value}</div>)}

      {hasLinks && (
        <DataLinksContextMenu links={() => getCellLinks(field, row) || []}>
          {(api) => {
            if (api.openMenu) {
              return (
                <button
                  className={cx(clearButtonStyle, getLinkStyle(tableStyles, cellOptions, api.targetClassName))}
                  onClick={api.openMenu}
                >
                  {value}
                </button>
              );
            } else {
              return <div className={getLinkStyle(tableStyles, cellOptions, api.targetClassName)}>{value}</div>;
            }
          }}
        </DataLinksContextMenu>
      )}

      {hover && showActions && <CellActions {...props} previewMode="text" showFilters={showFilters} />}
    </div>
  );
};

function getCellStyle(
  tableStyles: TableStyles,
  cellOptions: TableCellOptions,
  displayValue: DisplayValue,
  disableOverflowOnHover = false,
  isStringValue = false
) {
  // Setup color variables
  let textColor: string | undefined = undefined;
  let bgColor: string | undefined = undefined;

  // Get colors
  const colors = getCellColors(tableStyles, cellOptions, displayValue);
  textColor = colors.textColor;
  
  // Skip application of background color if we're applying to the whole row
  if (cellOptions.type === TableCellDisplayMode.ColorBackground && !cellOptions.applyToRow) {
    bgColor = colors.bgColor;
  }

  // If we have definied colors return those styles
  // Otherwise we return default styles
  if (textColor !== undefined || bgColor !== undefined) {
    return tableStyles.buildCellContainerStyle(textColor, bgColor, !disableOverflowOnHover, isStringValue);
  }

  if (isStringValue) {
    return disableOverflowOnHover ? tableStyles.cellContainerTextNoOverflow : tableStyles.cellContainerText;
  } else {
    return disableOverflowOnHover ? tableStyles.cellContainerNoOverflow : tableStyles.cellContainer;
  }
}

function getLinkStyle(tableStyles: TableStyles, cellOptions: TableCellOptions, targetClassName: string | undefined) {
  if (cellOptions.type === TableCellDisplayMode.Auto) {
    return cx(tableStyles.cellLink, targetClassName);
  }

  return cx(tableStyles.cellLinkForColoredCell, targetClassName);
}
