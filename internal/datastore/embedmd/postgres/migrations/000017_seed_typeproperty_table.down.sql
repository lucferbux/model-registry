-- Clear seeded TypeProperty data
DELETE FROM "TypeProperty"
WHERE (type_id, name, data_type) IN (
    ((SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'description', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'owner', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'state', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ModelVersion'), 'author', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ModelVersion'), 'description', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ModelVersion'), 'model_name', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ModelVersion'), 'state', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ModelVersion'), 'version', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.DocArtifact'), 'description', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ModelArtifact'), 'description', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ModelArtifact'), 'model_format_name', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ModelArtifact'), 'model_format_version', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ModelArtifact'), 'service_account_name', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ModelArtifact'), 'storage_key', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ModelArtifact'), 'storage_path', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ServingEnvironment'), 'description', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.InferenceService'), 'description', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.InferenceService'), 'desired_state', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.InferenceService'), 'model_version_id', 1),
    ((SELECT id FROM "Type" WHERE name = 'kf.InferenceService'), 'registered_model_id', 1),
    ((SELECT id FROM "Type" WHERE name = 'kf.InferenceService'), 'runtime', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.InferenceService'), 'serving_environment_id', 1),
    ((SELECT id FROM "Type" WHERE name = 'kf.ServeModel'), 'description', 3),
    ((SELECT id FROM "Type" WHERE name = 'kf.ServeModel'), 'model_version_id', 1)
); 