package cli

import (
	"fincli/internal/bankcsv/autofill"
	"os"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_defaultConfig_GetAutofillRules(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		config  string // YAML configuration string to be tested
		want    []autofill.Rule
		wantErr bool
	}{
		{
			name:   "no autofill key",
			config: "",
		},
		{
			name:   "autofill with no rules",
			config: "autofill:\n",
		},
		{
			name:   "autofill with empty rules",
			config: "autofill:\n  rules:\n",
		},
		{
			name: "autofill with one rule",
			config: heredoc.Doc(`
				autofill:
				  rules:
				    - name: Test Rule
				`),
			want: []autofill.Rule{
				{
					Name: "Test Rule",
				},
			},
		},
		{
			name: "autofill with multiple rules",
			config: heredoc.Doc(`
				autofill:
				  rules:
				    - name: Test Rule 1
				    - name: Test Rule 2
			`),
			want: []autofill.Rule{
				{
					Name: "Test Rule 1",
				},
				{
					Name: "Test Rule 2",
				},
			},
		},
		{
			name: "basic autofill rule",
			config: heredoc.Doc(`
				autofill:
				  rules:
				    - name: test rule
				      format: sample format 
				      target:
				        header: sample header title
				      value: test value
				      overwrite: never
				      condition:
				        field:
				          header: sample field header title
				        equals: sample equals value

				`),
			want: []autofill.Rule{
				{
					Name:   "test rule",
					Format: "sample format",
					Target: autofill.FieldRef{
						Header: "sample header title",
					},
					Value:     "test value",
					Overwrite: autofill.OverwriteNever,
					Condition: autofill.Condition{
						Field: autofill.FieldRef{
							Header: "sample field header title",
						},
						Type:  autofill.CondEquals,
						Value: "sample equals value",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfgFile := makeTestConfigFile(t, tt.config)
			d := New()
			d.SetConfigFile(cfgFile)
			got, gotErr := d.GetAutofillRules()
			require := require.New(t)
			assert := assert.New(t)
			if tt.wantErr {
				require.Error(gotErr)
				return
			}
			require.NoError(gotErr, "yaml was:\n%s", tt.config)
			require.Len(got, len(tt.want))
			for i, rule := range got {
				assert.Equal(tt.want[i].Name, rule.Name)
				assert.Equal(tt.want[i].Format, rule.Format)

				assert.Equal(tt.want[i].Target.Header, rule.Target.Header)
				assert.Equal(tt.want[i].Target.Position, rule.Target.Position)

				assert.Equal(tt.want[i].Value, rule.Value)
				assert.Equal(tt.want[i].Overwrite, rule.Overwrite)

				assert.Equal(tt.want[i].Condition.Field.Header, rule.Condition.Field.Header)
				assert.Equal(tt.want[i].Condition.Field.Position, rule.Condition.Field.Position)

				assert.Equal(tt.want[i].Condition.Type, rule.Condition.Type)
				assert.Equal(tt.want[i].Condition.Value, rule.Condition.Value)
			}
		})
	}
}

func makeTestConfigFile(t *testing.T, input string) (filename string) {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-config-*.yaml")
	require.NoError(t, err, "failed to create temp config file")
	defer f.Close()
	_, err = f.WriteString(input)
	require.NoError(t, err, "failed to write to temp config file")
	filename = f.Name()
	return filename
}
