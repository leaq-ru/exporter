package call

import (
	"github.com/nnqq/scr-proto/codegen/go/category"
	"github.com/nnqq/scr-proto/codegen/go/city"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"google.golang.org/grpc"
)

func NewClients(
	parserURL,
	cityURL,
	categoryURL string,
) (
	companyClient parser.CompanyClient,
	cityClient city.CityClient,
	categoryClient category.CategoryClient,
	err error,
) {
	connParser, err := grpc.Dial(parserURL, grpc.WithInsecure())
	if err != nil {
		return
	}
	companyClient = parser.NewCompanyClient(connParser)

	connCity, err := grpc.Dial(cityURL, grpc.WithInsecure())
	if err != nil {
		return
	}
	cityClient = city.NewCityClient(connCity)

	connCategory, err := grpc.Dial(categoryURL, grpc.WithInsecure())
	if err != nil {
		return
	}
	categoryClient = category.NewCategoryClient(connCategory)
	return
}
