package schema

import (
	"errors"
	"net/url"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/unbindapp/unbind-api/ent/schema/mixin"
)

// Registry holds the schema definition for the Registry entity.
type Registry struct {
	ent.Schema
}

// Mixin of the Registry.
func (Registry) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the Registry.
func (Registry) Fields() []ent.Field {
	return []ent.Field{
		field.String("host").Validate(func(s string) error {
			// Should be valid host without protocol
			// e.g. registry.example.com , registry.example.com:5000
			// but not https://registry.example.com
			parsedURL, err := url.Parse(s)
			if err != nil {
				return errors.New("invalid host format")
			}
			if parsedURL.Scheme != "" {
				return errors.New("host should not contain protocol")
			}
			if parsedURL.Hostname() == "" {
				return errors.New("host should not be empty")
			}
			return nil
		}),
		field.String("kubernetes_secret").Optional().Nillable().
			Comment(
				"The name of the kubernetes registry credentials secret, should be located in the unbind system namespace",
			),
		field.Bool("is_default").
			Comment(
				"If true, this is the registry that will be used for internal CI/CD",
			),
	}
}

// Edges of the Registry.
func (Registry) Edges() []ent.Edge {
	return nil
}

// Annotations of the Registry
func (Registry) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "registries",
		},
	}
}
