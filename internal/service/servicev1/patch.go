package servicev1

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"

	"go.verv.tech/matreshka-be/internal/domain"
	"go.verv.tech/matreshka-be/internal/service/user_errors"
)

const (
	appInfo            = "APP-INFO"
	environmentSegment = "ENVIRONMENT"
	dataSourceSegment  = "SERVERS"
	serverSegment      = "DATA-SOURCES"
)

func (c *ConfigService) PatchConfig(ctx context.Context, req domain.PatchConfigRequest) error {
	patchToUpdate := newPatch(req.Batch)

	cfgNodes, err := c.configStorage.GetConfigNodes(ctx, req.ServiceName)
	if err != nil {
		return errors.Wrap(err, "error getting nodes")
	}

	if cfgNodes == nil {
		_, err = c.CreateConfig(ctx, req.ServiceName)
		if err != nil {
			return errors.Wrap(err, "error creating config to patch to")
		}
		cfgNodes = &evon.Node{}
	}

	patchToUpdate.normalizeEnvironmentChanges(cfgNodes)

	err = c.txManager.Execute(func(tx *sql.Tx) error {
		configStorage := c.configStorage.WithTx(tx)

		err = configStorage.DeleteValues(ctx, req.ServiceName, patchToUpdate.delete)
		if err != nil {
			return errors.Wrap(err, "error deleting values")
		}

		upsertReq := append(patchToUpdate.upsert, patchToUpdate.envUpsert...)
		err = configStorage.UpsertValues(ctx, req.ServiceName, upsertReq)
		if err != nil {
			return errors.Wrap(err, "error patching config in data storage")
		}

		err = configStorage.SetUpdatedAt(ctx, req.ServiceName, time.Now())
		if err != nil {
			return errors.Wrap(err, "error updating time")
		}

		return nil
	})

	if len(patchToUpdate.invalid) != 0 {
		return errors.Wrap(user_errors.ErrValidation, "Invalid patched env var name: "+fmt.Sprint(patchToUpdate.invalid))
	}

	return nil
}

func (p *patch) normalizeEnvironmentChanges(cfg *evon.Node) {
	nodeStorage := evon.NodesToStorage(cfg.InnerNodes)

	newEnvValues := make(map[string]domain.PatchConfig)
	typesMap := make(map[string]domain.PatchConfig)
	enumMap := make(map[string]domain.PatchConfig)

	envUpsert := make([]domain.PatchConfig, 0, len(newEnvValues))

	for _, valuePatch := range p.envUpsert {
		// already exists -> simply update value
		_, ok := nodeStorage[valuePatch.FieldName]
		if ok {
			envUpsert = append(envUpsert, valuePatch)
			continue
		}

		if strings.HasSuffix(valuePatch.FieldName, "_TYPE") {
			typesMap[valuePatch.FieldName[:len(valuePatch.FieldName)-5]] = valuePatch
			continue
		}

		if strings.HasSuffix(valuePatch.FieldName, "_ENUM") {
			enumMap[valuePatch.FieldName] = valuePatch
			continue
		}

		newEnvValues[valuePatch.FieldName] = valuePatch
	}

	for key, patchVal := range newEnvValues {
		typeVal, ok := typesMap[key]
		if !ok {
			p.invalid = append(p.invalid, patchVal)
			continue
		}

		envUpsert = append(envUpsert, patchVal, typeVal)
		enumVal, ok := enumMap[key]
		if ok {
			envUpsert = append(envUpsert, enumVal)
		}
	}

	p.envUpsert = envUpsert
}

type patch struct {
	invalid   []domain.PatchConfig
	upsert    []domain.PatchConfig
	envUpsert []domain.PatchConfig
	delete    []domain.PatchConfig
}

func newPatch(batch []domain.PatchConfig) patch {
	p := patch{}
	for _, ptch := range batch {
		var ok bool
		ptch.FieldName, ok = validateName(ptch)
		if !ok {
			p.invalid = append(p.invalid, ptch)
			continue
		}

		val := extractValue(ptch.FieldValue)
		if val == nil {
			p.delete = append(p.delete, ptch)
		} else {
			if strings.HasPrefix(ptch.FieldName, environmentSegment) {
				p.envUpsert = append(p.envUpsert, ptch)
			} else {
				p.upsert = append(p.upsert, ptch)
			}

		}
	}
	return p
}

func extractValue(in any) any {
	inRef := reflect.ValueOf(in)
	if inRef.IsNil() {
		return nil
	}

	if inRef.Kind() != reflect.Ptr {
		return inRef.Interface()
	}

	return inRef.Elem().Interface()
}

func validateName(patch domain.PatchConfig) (newName string, ok bool) {
	for _, segment := range allowedSegments {
		if strings.HasPrefix(patch.FieldName, segment) {
			return strings.ToUpper(patch.FieldName), true
		}
	}

	return patch.FieldName, false
}

func clearName(name string) string {
	fromIdx := strings.Index(name, "/")
	if fromIdx == -1 {
		return name
	}
	newName := (name)[fromIdx+1:]
	newName = strings.ReplaceAll(newName, "/", "-")
	return newName
}
