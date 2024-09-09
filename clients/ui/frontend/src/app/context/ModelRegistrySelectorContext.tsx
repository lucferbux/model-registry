import * as React from 'react';
import { ModelRegistry } from '~/app/types';
import useModelRegistries from '~/app/hooks/useModelRegistries';

export type ModelRegistrySelectorContextType = {
  modelRegistriesLoaded: boolean;
  modelRegistriesLoadError?: Error;
  modelRegistries: ModelRegistry[];
  preferredModelRegistry: ModelRegistry | undefined;
  updatePreferredModelRegistry: (modelRegistry: ModelRegistry | undefined) => void;
};

type ModelRegistrySelectorContextProviderProps = {
  children: React.ReactNode;
};

export const ModelRegistrySelectorContext = React.createContext<ModelRegistrySelectorContextType>({
  modelRegistriesLoaded: false,
  modelRegistriesLoadError: undefined,
  modelRegistries: [],
  preferredModelRegistry: undefined,
  updatePreferredModelRegistry: () => undefined,
});

export const ModelRegistrySelectorContextProvider: React.FC<
  ModelRegistrySelectorContextProviderProps
> = ({ children, ...props }) => (
  <EnabledModelRegistrySelectorContextProvider {...props}>
    {children}
  </EnabledModelRegistrySelectorContextProvider>
);

const EnabledModelRegistrySelectorContextProvider: React.FC<
  ModelRegistrySelectorContextProviderProps
> = ({ children }) => {
  const [modelRegistries, isLoaded, error] = useModelRegistries();
  const [preferredModelRegistry, setPreferredModelRegistry] =
    React.useState<ModelRegistrySelectorContextType['preferredModelRegistry']>(undefined);

  const firstModelRegistry = modelRegistries.length > 0 ? modelRegistries[0] : null;

  return (
    <ModelRegistrySelectorContext.Provider
      value={React.useMemo(
        () => ({
          modelRegistriesLoaded: isLoaded,
          modelRegistriesLoadError: error,
          modelRegistries,
          preferredModelRegistry: preferredModelRegistry ?? firstModelRegistry ?? undefined,
          updatePreferredModelRegistry: setPreferredModelRegistry,
        }),
        [isLoaded, error, modelRegistries, preferredModelRegistry, firstModelRegistry],
      )}
    >
      {children}
    </ModelRegistrySelectorContext.Provider>
  );
};
