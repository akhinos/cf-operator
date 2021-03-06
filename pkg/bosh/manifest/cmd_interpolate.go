package manifest

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/cppforlife/go-patch/patch"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// InterpolateVariables reads explicit secrets from a folder and writes an interpolated manifest to the output.json file in /mnt/quarks volume mount.
func InterpolateVariables(log *zap.SugaredLogger, boshManifestBytes []byte, variablesDir string, outputFilePath string) error {
	var vars []boshtpl.Variables

	variables, err := ioutil.ReadDir(variablesDir)
	if err != nil {
		return errors.Wrapf(err, "could not read variables directory")
	}

	for _, variable := range variables {
		// Each directory is a variable name
		if variable.IsDir() {
			staticVars := boshtpl.StaticVariables{}
			// Each filename is a field name and its context is a variable value
			err = filepath.Walk(filepath.Clean(variablesDir+"/"+variable.Name()), func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() {
					_, varFileName := filepath.Split(path)
					// Skip the symlink to a directory
					if strings.HasPrefix(varFileName, "..") {
						return filepath.SkipDir
					}
					varBytes, err := ioutil.ReadFile(path)
					if err != nil {
						return errors.Wrapf(err, "could not read variables variable %s", variable.Name())
					}

					// If variable type is password, set password value directly
					switch varFileName {
					case "password":
						staticVars[variable.Name()] = string(varBytes)
					default:
						staticVars[variable.Name()] = mergeStaticVar(staticVars[variable.Name()], varFileName, string(varBytes))
					}
				}
				return nil
			})
			if err != nil {
				return errors.Wrapf(err, "could not read directory  %s", variable.Name())
			}

			vars = append(vars, staticVars)
		}
	}

	multiVars := boshtpl.NewMultiVars(vars)
	tpl := boshtpl.NewTemplate(boshManifestBytes)

	// Following options are empty for cf-operator
	op := patch.Ops{}
	evalOpts := boshtpl.EvaluateOpts{
		ExpectAllKeys:     false,
		ExpectAllVarsUsed: false,
	}

	yamlBytes, err := tpl.Evaluate(multiVars, op, evalOpts)
	if err != nil {
		return errors.Wrapf(err, "could not evaluate variables")
	}

	m, err := LoadYAML(yamlBytes)
	if err != nil {
		return errors.Wrapf(err, "could not evaluate variables")
	}

	yamlBytes, err = m.Marshal()
	if err != nil {
		return errors.Wrapf(err, "could not evaluate variables")
	}

	jsonBytes, err := json.Marshal(map[string]string{
		DesiredManifestKeyName: string(yamlBytes),
	})
	if err != nil {
		return errors.Wrapf(err, "could not marshal json output")
	}

	err = ioutil.WriteFile(outputFilePath, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func mergeStaticVar(staticVar interface{}, field string, value string) interface{} {
	if staticVar == nil {
		staticVar = map[interface{}]interface{}{
			field: value,
		}
	} else {
		staticVarMap := staticVar.(map[interface{}]interface{})
		staticVarMap[field] = value
		staticVar = staticVarMap
	}

	return staticVar
}
