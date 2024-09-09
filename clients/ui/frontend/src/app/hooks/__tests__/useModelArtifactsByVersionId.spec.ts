import useModelArtifactsByVersionId from '~/app/hooks/useModelArtifactsByVersionId';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';
import { NotReadyError } from '~/utilities/useFetchState';
import { ModelRegistryAPIs } from '~/app/types';
import { mockModelArtifact } from '~/__mocks__/mockModelArtifact';
import { testHook } from '~/__tests__/unit/testUtils/hooks';

global.fetch = jest.fn();
// Mock the useModelRegistryAPI hook
jest.mock('~/app/hooks/useModelRegistryAPI', () => ({
  useModelRegistryAPI: jest.fn(),
}));

const mockUseModelRegistryAPI = jest.mocked(useModelRegistryAPI);

const mockModelRegistryAPIs: ModelRegistryAPIs = {
  createRegisteredModel: jest.fn(),
  createModelVersionForRegisteredModel: jest.fn(),
  createModelArtifactForModelVersion: jest.fn(),
  getRegisteredModel: jest.fn(),
  getModelVersion: jest.fn(),
  listRegisteredModels: jest.fn(),
  getModelVersionsByRegisteredModel: jest.fn(),
  getModelArtifactsByModelVersion: jest.fn(),
  patchRegisteredModel: jest.fn(),
  patchModelVersion: jest.fn(),
};

describe('useModelArtifactsByVersionId', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should return NotReadyError if API is not available', async () => {
    mockUseModelRegistryAPI.mockReturnValue({
      api: mockModelRegistryAPIs,
      apiAvailable: false,
      refreshAllAPI: jest.fn(),
    });

    const { result } = testHook(useModelArtifactsByVersionId)();
    const [, , error] = result.current;

    expect(error).toBe('API not yet available');
    expect(error?.message).toBeInstanceOf(NotReadyError);
  });

  it('should return NotReadyError if modelVersionId is not provided', async () => {
    mockUseModelRegistryAPI.mockReturnValue({
      api: mockModelRegistryAPIs,
      apiAvailable: true,
      refreshAllAPI: jest.fn(),
    });

    const { result } = testHook(useModelArtifactsByVersionId)();
    const [, , error] = result.current;

    expect(error).toBeInstanceOf(NotReadyError);
    expect(error?.message).toBe('No model registeredModel id');
  });

  it('should fetch model artifacts if API is available and modelVersionId is provided', async () => {
    const mockedResponse = {
      items: [mockModelArtifact({ id: 'artifact-1' })],
      size: 1,
      pageSize: 1,
    };
    const mockGetModelArtifactsByModelVersion = jest.fn().mockResolvedValue(mockedResponse);

    mockUseModelRegistryAPI.mockReturnValue({
      api: {
        ...mockModelRegistryAPIs,
        getModelArtifactsByModelVersion: mockGetModelArtifactsByModelVersion,
      },
      apiAvailable: false,
      refreshAllAPI: jest.fn(),
    });

    const { result } = testHook(useModelArtifactsByVersionId)('version-id');
    const [data] = result.current;

    expect(data).toEqual(mockedResponse);
    expect(mockGetModelArtifactsByModelVersion).toHaveBeenCalledWith(
      expect.any(Object),
      'version-id',
    );
  });
});
