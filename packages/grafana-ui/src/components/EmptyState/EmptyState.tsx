import React, { ReactNode } from 'react';

import { t } from '../../utils/i18n';
import { Button } from '../Button';
import { Box } from '../Layout/Box/Box';
import { Stack } from '../Layout/Stack/Stack';
import { Text } from '../Text/Text';

import { GrotNotFound } from './GrotNotFound/GrotNotFound';

interface BaseProps {
  /**
   * Provide a button to render below the message
   */
  button?: ReactNode;
  hideImage?: boolean;
  /**
   * Override the default image for the variant
   */
  image?: ReactNode;
  /**
   * Message to display to the user
   */
  message?: string;
  /**
   * Empty state variant. Possible values are 'search'.
   */
  variant?: 'default' | 'search';
}

interface DefaultVariantProps extends BaseProps {
  message: string;
  variant?: 'default';
}

interface SearchVariantProps extends BaseProps {
  variant: 'search';
}

type Props = DefaultVariantProps | SearchVariantProps;

export const EmptyState = ({
  button,
  children,
  image,
  message = t('grafana-ui.empty-state.search-message', 'No results found'),
  hideImage = false,
  variant = 'default',
}: React.PropsWithChildren<Props>) => {
  const imageToShow = image ?? getImageForVariant(variant);

  return (
    <Box paddingY={4} gap={4} display="flex" direction="column" alignItems="center">
      {!hideImage && imageToShow}
      <Stack direction="column" alignItems="center">
        <Text variant="h4">{message}</Text>
        {children && <Text color="secondary">{children}</Text>}
      </Stack>
      {button}
    </Box>
  );
};

function getImageForVariant(variant: Props['variant']) {
  if (variant === 'search') {
    return <GrotNotFound width={300} />;
  }

  return <Button>foo</Button>;
}