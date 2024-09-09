package core

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// ARTIFACTS

// UpsertArtifact creates a new artifact if the provided artifact's ID is nil, or updates an existing artifact if the
// ID is provided.
// A model version ID must be provided to disambiguate between artifacts.
// Upon creation, new artifacts will be associated with their corresponding model version.
func (serv *ModelRegistryService) UpsertArtifact(artifact *openapi.Artifact, modelVersionId *string) (*openapi.Artifact, error) {
	if artifact == nil {
		return nil, fmt.Errorf("invalid artifact pointer, can't upsert nil")
	}
	creating := false
	if ma := artifact.ModelArtifact; ma != nil {
		if ma.Id == nil {
			creating = true
			glog.Info("Creating model artifact")
			if modelVersionId == nil {
				return nil, fmt.Errorf("missing model version id, cannot create artifact without model version: %w", api.ErrBadRequest)
			}
			_, err := serv.GetModelVersionById(*modelVersionId)
			if err != nil {
				return nil, fmt.Errorf("no model version found for id %s: %w", *modelVersionId, api.ErrNotFound)
			}
		} else {
			glog.Info("Updating model artifact")
			existing, err := serv.GetModelArtifactById(*ma.Id)
			if err != nil {
				return nil, err
			}

			withNotEditable, err := serv.openapiConv.OverrideNotEditableForModelArtifact(converter.NewOpenapiUpdateWrapper(existing, ma))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}
			ma = &withNotEditable

			_, err = serv.getModelVersionByArtifactId(*ma.Id)
			if err != nil {
				return nil, err
			}
		}
	} else if da := artifact.DocArtifact; da != nil {
		if da.Id == nil {
			creating = true
			glog.Info("Creating doc artifact")
			if modelVersionId == nil {
				return nil, fmt.Errorf("missing model version id, cannot create artifact without model version: %w", api.ErrBadRequest)
			}
			_, err := serv.GetModelVersionById(*modelVersionId)
			if err != nil {
				return nil, fmt.Errorf("no model version found for id %s: %w", *modelVersionId, api.ErrNotFound)
			}
		} else {
			glog.Info("Updating doc artifact")
			existing, err := serv.GetArtifactById(*da.Id)
			if err != nil {
				return nil, err
			}
			if existing.DocArtifact == nil {
				return nil, fmt.Errorf("mismatched types, artifact with id %s is not a doc artifact: %w", *da.Id, api.ErrBadRequest)
			}

			withNotEditable, err := serv.openapiConv.OverrideNotEditableForDocArtifact(converter.NewOpenapiUpdateWrapper(existing.DocArtifact, da))
			if err != nil {
				return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
			}
			da = &withNotEditable

			_, err = serv.getModelVersionByArtifactId(*da.Id)
			if err != nil {
				return nil, err
			}
		}
	} else {
		return nil, fmt.Errorf("invalid artifact type, must be either ModelArtifact or DocArtifact: %w", api.ErrBadRequest)
	}
	pa, err := serv.mapper.MapFromArtifact(artifact, modelVersionId)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	artifactsResp, err := serv.mlmdClient.PutArtifacts(context.Background(), &proto.PutArtifactsRequest{
		Artifacts: []*proto.Artifact{pa},
	})
	if err != nil {
		return nil, err
	}

	if creating {
		// add explicit Attribution between Artifact and ModelVersion
		modelVersionId, err := converter.StringToInt64(modelVersionId)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		attributions := []*proto.Attribution{}
		for _, a := range artifactsResp.ArtifactIds {
			attributions = append(attributions, &proto.Attribution{
				ContextId:  modelVersionId,
				ArtifactId: &a,
			})
		}
		_, err = serv.mlmdClient.PutAttributionsAndAssociations(context.Background(), &proto.PutAttributionsAndAssociationsRequest{
			Attributions: attributions,
			Associations: make([]*proto.Association, 0),
		})
		if err != nil {
			return nil, err
		}
	}

	idAsString := converter.Int64ToString(&artifactsResp.ArtifactIds[0])
	return serv.GetArtifactById(*idAsString)
}

func (serv *ModelRegistryService) GetArtifactById(id string) (*openapi.Artifact, error) {
	idAsInt, err := converter.StringToInt64(&id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	artifactsResp, err := serv.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{int64(*idAsInt)},
	})
	if err != nil {
		return nil, err
	}
	if len(artifactsResp.Artifacts) > 1 {
		return nil, fmt.Errorf("multiple artifacts found for id %s: %w", id, api.ErrNotFound)
	}
	if len(artifactsResp.Artifacts) == 0 {
		return nil, fmt.Errorf("no artifact found for id %s: %w", id, api.ErrNotFound)
	}
	return serv.mapper.MapToArtifact(artifactsResp.Artifacts[0])
}

func (serv *ModelRegistryService) GetArtifacts(listOptions api.ListOptions, modelVersionId *string) (*openapi.ArtifactList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	var artifacts []*proto.Artifact
	var nextPageToken *string
	if modelVersionId == nil {
		return nil, fmt.Errorf("missing model version id, cannot get artifacts without model version: %w", api.ErrBadRequest)
	}
	ctxId, err := converter.StringToInt64(modelVersionId)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}
	artifactsResp, err := serv.mlmdClient.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
		ContextId: ctxId,
		Options:   listOperationOptions,
	})
	if err != nil {
		return nil, err
	}
	artifacts = artifactsResp.Artifacts
	nextPageToken = artifactsResp.NextPageToken

	results := []openapi.Artifact{}
	for _, a := range artifacts {
		mapped, err := serv.mapper.MapToArtifact(a)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ArtifactList{
		NextPageToken: apiutils.ZeroIfNil(nextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}

// MODEL ARTIFACTS

// UpsertModelArtifact creates a new model artifact if the provided model artifact's ID is nil,
// or updates an existing model artifact if the ID is provided.
// If a model version ID is provided and the model artifact is newly created, establishes an
// explicit attribution between the model version and the created model artifact.
func (serv *ModelRegistryService) UpsertModelArtifact(modelArtifact *openapi.ModelArtifact, modelVersionId *string) (*openapi.ModelArtifact, error) {
	art, err := serv.UpsertArtifact(&openapi.Artifact{
		ModelArtifact: modelArtifact,
	}, modelVersionId)
	if err != nil {
		return nil, err
	}
	return art.ModelArtifact, err
}

// GetModelArtifactById retrieves a model artifact by its unique identifier (ID).
func (serv *ModelRegistryService) GetModelArtifactById(id string) (*openapi.ModelArtifact, error) {
	art, err := serv.GetArtifactById(id)
	if err != nil {
		return nil, err
	}
	ma := art.ModelArtifact
	if ma == nil {
		return nil, fmt.Errorf("artifact with id %s is not a model artifact: %w", id, api.ErrNotFound)
	}
	return ma, err
}

// GetModelArtifactByInferenceService retrieves the model artifact associated with the specified inference service ID.
func (serv *ModelRegistryService) GetModelArtifactByInferenceService(inferenceServiceId string) (*openapi.ModelArtifact, error) {
	mv, err := serv.GetModelVersionByInferenceService(inferenceServiceId)
	if err != nil {
		return nil, err
	}

	artifactList, err := serv.GetModelArtifacts(api.ListOptions{}, mv.Id)
	if err != nil {
		return nil, err
	}

	if artifactList.Size == 0 {
		return nil, fmt.Errorf("no artifacts found for model version %s: %w", *mv.Id, api.ErrNotFound)
	}

	return &artifactList.Items[0], nil
}

// GetModelArtifactByParams retrieves a model artifact based on specified parameters, such as (artifact name and model version ID), or external ID.
// If multiple or no model artifacts are found, an error is returned.
func (serv *ModelRegistryService) GetModelArtifactByParams(artifactName *string, modelVersionId *string, externalId *string) (*openapi.ModelArtifact, error) {
	var artifact0 *proto.Artifact

	filterQuery := ""
	if externalId != nil {
		filterQuery = fmt.Sprintf("external_id = \"%s\"", *externalId)
	} else if artifactName != nil && modelVersionId != nil {
		filterQuery = fmt.Sprintf("name = \"%s\"", converter.PrefixWhenOwned(modelVersionId, *artifactName))
	} else {
		return nil, fmt.Errorf("invalid parameters call, supply either (artifactName and modelVersionId), or externalId: %w", api.ErrBadRequest)
	}
	glog.Info("filterQuery ", filterQuery)

	artifactsResponse, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
		TypeName: &serv.nameConfig.ModelArtifactTypeName,
		Options: &proto.ListOperationOptions{
			FilterQuery: &filterQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(artifactsResponse.Artifacts) > 1 {
		return nil, fmt.Errorf("multiple model artifacts found for artifactName=%v, modelVersionId=%v, externalId=%v: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(modelVersionId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	if len(artifactsResponse.Artifacts) == 0 {
		return nil, fmt.Errorf("no model artifacts found for artifactName=%v, modelVersionId=%v, externalId=%v: %w", apiutils.ZeroIfNil(artifactName), apiutils.ZeroIfNil(modelVersionId), apiutils.ZeroIfNil(externalId), api.ErrNotFound)
	}

	artifact0 = artifactsResponse.Artifacts[0]

	result, err := serv.mapper.MapToModelArtifact(artifact0)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	return result, nil
}

// GetModelArtifacts retrieves a list of model artifacts based on the provided list options and optional model version ID.
func (serv *ModelRegistryService) GetModelArtifacts(listOptions api.ListOptions, modelVersionId *string) (*openapi.ModelArtifactList, error) {
	listOperationOptions, err := apiutils.BuildListOperationOptions(listOptions)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
	}

	var artifacts []*proto.Artifact
	var nextPageToken *string
	if modelVersionId != nil {
		ctxId, err := converter.StringToInt64(modelVersionId)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		artifactsResp, err := serv.mlmdClient.GetArtifactsByContext(context.Background(), &proto.GetArtifactsByContextRequest{
			ContextId: ctxId,
			Options:   listOperationOptions,
		})
		if err != nil {
			return nil, err
		}
		artifacts = artifactsResp.Artifacts
		nextPageToken = artifactsResp.NextPageToken
	} else {
		artifactsResp, err := serv.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
			TypeName: &serv.nameConfig.ModelArtifactTypeName,
			Options:  listOperationOptions,
		})
		if err != nil {
			return nil, err
		}
		artifacts = artifactsResp.Artifacts
		nextPageToken = artifactsResp.NextPageToken
	}

	results := []openapi.ModelArtifact{}
	for _, a := range artifacts {
		mapped, err := serv.mapper.MapToModelArtifact(a)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		results = append(results, *mapped)
	}

	toReturn := openapi.ModelArtifactList{
		NextPageToken: apiutils.ZeroIfNil(nextPageToken),
		PageSize:      apiutils.ZeroIfNil(listOptions.PageSize),
		Size:          int32(len(results)),
		Items:         results,
	}
	return &toReturn, nil
}