package incidentio

import (
	"regexp"

	"github.com/ricoberger/grafana-incidentio-plugin/pkg/models"
)

func getAlertAttribute(attribute string, alertAttributes []models.AlertAttribute) *models.AlertAttribute {
	for _, attr := range alertAttributes {
		if attr.ID == attribute {
			return &attr
		}
	}
	return nil
}

func getCatalogEntryFromAlertAttribute(alertAttribute *models.AlertAttribute) string {
	re := regexp.MustCompile(`CatalogEntry\["([^"]+)"\]`)
	catalogEntryMatches := re.FindStringSubmatch(alertAttribute.Type)

	if len(catalogEntryMatches) > 1 {
		return catalogEntryMatches[1]
	}
	return ""
}

func getCatalogType(alertAttribute *models.AlertAttribute, catalogTypes []models.CatalogType) *models.CatalogType {
	catalogEntry := getCatalogEntryFromAlertAttribute(alertAttribute)
	if catalogEntry == "" {
		return nil
	}

	for _, ct := range catalogTypes {
		if ct.ID == catalogEntry {
			return &ct
		} else if ct.TypeName == catalogEntry {
			return &ct
		}
	}
	return nil
}

func getCustomField(customFieldID string, customFields []models.CustomField) *models.CustomField {
	for _, cf := range customFields {
		if cf.ID == customFieldID {
			return &cf
		}
	}
	return nil
}
