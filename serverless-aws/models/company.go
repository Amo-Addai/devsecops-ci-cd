package models

import (
	"log"

	"github.com/Amo-Addai/devsecops-ci-cd/serverless-aws/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/guregu/dynamo"
)

// Company - the dynamodb representation of a company
type Company struct {
	Phone      string `dynamo:"phone"`
	Name       string `dynamo:"name"`
	CompanyID  string `dynamo:"company_id"`
	FilterName string `dynamo:"filter_name"`
	BotURL     string `dynamo:"bot_url"`
	OpenTime   string `dynamo:"open_time,omitempty"`
	CloseTime  string `dynamo:"close_time,omitempty"`
	BotEnabled bool   `dynamo:"bot_enabled"`
}

// CompanyJSON - json representation of a company
type CompanyJSON struct {
	Phone      string `json:"phone"`
	Name       string `json:"name"`
	CompanyID  string `json:"company_id"`
	FilterName string `json:"filter_name"`
	BotURL     string `json:"bot_url"`
	OpenTime   string `json:"open_time,omitempty"`
	CloseTime  string `json:"close_time,omitempty"`
	BotEnabled bool   `json:"bot_enabled"`
}

// ConvertToCompanyJSON - convert the conversationMessage to a json representable object
func (company *Company) ConvertToCompanyJSON() *CompanyJSON {
	return &CompanyJSON{
		CompanyID:  company.CompanyID,
		Name:       company.Name,
		FilterName: company.FilterName,
		BotURL:     company.BotURL,
		Phone:      company.Phone,
		OpenTime:   company.OpenTime,
		CloseTime:  company.CloseTime,
		BotEnabled: company.BotEnabled,
	}
}

// FindCompanyByPhone - Find company by it's unique phone
func FindCompanyByPhone(table dynamo.Table, phone string, debug bool) (*Company, error) {
	var errToReturn error
	var company *Company

	if debug {
		log.Printf("[FindCompanyByPhone] phone: %s", phone)
	}
	err := table.Get("phone", phone).Index("company_phone_index").One(&company)
	if err != nil {
		if err == dynamo.ErrNotFound {
			if debug {
				log.Print("[FindCompanyByPhone] No company found, from orm")
			}
			errToReturn = nil
		} else {
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
					if debug {
						log.Print("[FindCompanyByPhone] No company found, from dynamodb")
					}
					errToReturn = nil
				} else {
					if debug {
						log.Print("[FindCompanyByPhone] Bad aerr")
						log.Print(aerr)
					}
					errToReturn = aerr
				}
			} else {
				log.Print("[FindCompanyByPhone] Error while trying to find existing company")
				log.Print(err)
				errToReturn = err
			}
		}
		return company, errToReturn
	}
	if debug {
		log.Print("[FindCompanyByPhone] Company found!")
	}
	return company, nil
}

// OpenTimeHourMin - Gets the open time hour and minute
func (company *Company) OpenTimeHourMin() (hour int, min int) {
	return utils.ParseTimeFieldForHourMin(company.OpenTime)
}

// CloseTimeHourMin - Gets the close time hour and minute
func (company *Company) CloseTimeHourMin() (hour int, min int) {
	return utils.ParseTimeFieldForHourMin(company.CloseTime)
}

// GetCompanies - Gets the companies for the given company phones
func GetCompanies(debug bool, region string, tableName string, companyPhones []string) ([]Company, error) {
	var companies []Company
	if companyPhones == nil || len(companyPhones) == 0 {
		return companies, nil
	}

	// setup
	awscfg := &aws.Config{}
	awscfg.WithRegion(region)
	sess := session.Must(session.NewSession(awscfg))
	svc := dynamodb.New(sess)

	// build expression
	filt := expression.Name("phone").Equal(expression.Value(companyPhones[0]))
	companyPhonesLength := len(companyPhones)
	if companyPhonesLength > 1 {
		for i := 1; i < companyPhonesLength; i++ {
			phone := companyPhones[i]
			filt = filt.Or(expression.Name("phone").Equal(expression.Value(phone)))
		}
	}

	// build
	expr, buildExpressionErr := expression.NewBuilder().WithFilter(filt).Build()
	if buildExpressionErr != nil {
		log.Println("Error while building expression for retrieving companies by phone")
		return companies, buildExpressionErr
	}

	// create input
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 &tableName,
	}

	// scan
	result, err := svc.Scan(params)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return companies, nil
		}
		log.Println("Error while retrieving companies for company phones")
		return companies, err
	}

	// parse result
	items := make([]Company, *result.Count)
	for i, item := range result.Items {
		items[i] = Company{
			Phone:      utils.ExtractDynamoResultString(item, "phone"),
			Name:       utils.ExtractDynamoResultString(item, "name"),
			CompanyID:  utils.ExtractDynamoResultString(item, "company_id"),
			FilterName: utils.ExtractDynamoResultString(item, "filter_name"),
			BotURL:     utils.ExtractDynamoResultString(item, "bot_url"),
			OpenTime:   utils.ExtractDynamoResultString(item, "open_time"),
			CloseTime:  utils.ExtractDynamoResultString(item, "close_time"),
			BotEnabled: utils.ExtractDynamoResultBool(item, "bot_enabled"),
		}
	}
	return items, nil
}

// GetAllCompanies - Get all the companies in the system.  No paginagtion, be very careful with this method as it's unbounded.
func GetAllCompanies(debug bool, region string, tableName string) ([]Company, error) {
	var companies []Company
	var errToReturn error

	table := openDBConn(region, tableName)
	err := table.Scan().All(&companies)

	if err != nil {
		if err == dynamo.ErrNotFound {
			if debug {
				log.Print("[GetAllCompanies] No companies found, from orm")
			}
			errToReturn = nil
		} else {
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
					if debug {
						log.Print("[GetAllCompanies] No companies found, from dynamodb")
					}
					errToReturn = nil
				} else {
					if debug {
						log.Print("[GetAllCompanies] Bad aerr")
						log.Print(aerr)
					}
					errToReturn = aerr
				}
			} else {
				log.Print("[GetAllCompanies] Error while trying to find companies")
				log.Print(err)
				errToReturn = err
			}
		}
		return companies, errToReturn
	}
	if debug {
		log.Print("[GetAllCompanies] Companies found!")
	}
	return companies, nil
}

// Private

func openDBConn(region string, tableName string) dynamo.Table {
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String(region)})
	return db.Table(tableName)
}
