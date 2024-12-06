package resources

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/accessanalyzer"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AccessAnalyzerResource = "AccessAnalyzer"

func init() {
	registry.Register(&registry.Registration{
		Name:                AccessAnalyzerResource,
		Scope:               nuke.Account,
		Resource:            &AccessAnalyzer{},
		Lister:              &AccessAnalyzerLister{},
		AlternativeResource: "AWS::AccessAnalyzer::Analyzer",
	})
}

type AccessAnalyzerLister struct{}

func (l *AccessAnalyzerLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := accessanalyzer.New(opts.Session)

	params := &accessanalyzer.ListAnalyzersInput{
		Type: aws.String("ACCOUNT"),
	}

	resources := make([]resource.Resource, 0)
	if err := svc.ListAnalyzersPages(params,
		func(page *accessanalyzer.ListAnalyzersOutput, lastPage bool) bool {
			for _, analyzer := range page.Analyzers {
				resources = append(resources, &AccessAnalyzer{
					svc:    svc,
					ARN:    analyzer.Arn,
					Name:   analyzer.Name,
					Status: analyzer.Status,
					Tags:   analyzer.Tags,
				})
			}
			return true
		}); err != nil {
		return nil, err
	}

	return resources, nil
}

type AccessAnalyzer struct {
	svc    *accessanalyzer.AccessAnalyzer
	ARN    *string            `description:"The ARN of the analyzer"`
	Name   *string            `description:"The name of the analyzer"`
	Status *string            `description:"The status of the analyzer"`
	Tags   map[string]*string `description:"The tags of the analyzer"`
}

func (r *AccessAnalyzer) Remove(_ context.Context) error {
	_, err := r.svc.DeleteAnalyzer(&accessanalyzer.DeleteAnalyzerInput{AnalyzerName: r.Name})

	return err
}

func (r *AccessAnalyzer) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AccessAnalyzer) String() string {
	return *r.Name
}