import { ValueOf } from '~/shared/typeHelpers';

//  TODO: [Data Flow] Get the status config params
export type UserSettings = {
  userId: string;
  clusterAdmin?: boolean;
};

// TODO: [Data Flow] Add more config parameters
export type ConfigSettings = {
  common: CommonConfig;
};

// TODO: [Data Flow] Add more config parameters
export type CommonConfig = {
  featureFlags: FeatureFlag;
};

// TODO: [Data Flow] Add more config parameters
export type FeatureFlag = {
  modelRegistry: boolean;
};

export type KeyValuePair = {
  key: string;
  value: string;
};

export type UpdateObjectAtPropAndValue<T> = (propKey: keyof T, propValue: ValueOf<T>) => void;
