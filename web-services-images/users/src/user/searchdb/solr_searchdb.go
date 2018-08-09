package searchdb

import (
	"fmt"
	"web-service-users/src/model"
	"web-service-users/src/user"

	"github.com/vanng822/go-solr/solr"
)

type solrSearchDB struct {
	Interface *solr.SolrInterface
}

func NewSolrSearchDB(solrInterface *solr.SolrInterface) user.SearchDB {
	return &solrSearchDB{solrInterface}
}

func (s *solrSearchDB) docToUserFound(doc *solr.Document) *model.UserFound {
	return &model.UserFound{
		UserName: decodeArrayString(doc.Get("user_name").([]interface{})),
		Name:     decodeArrayString(doc.Get("name").([]interface{})),
	}
}

func decodeArrayString(arrayinterface []interface{}) []string {
	s := make([]string, len(arrayinterface))
	for i, v := range arrayinterface {
		s[i] = fmt.Sprint(v)
	}
	return s
}

func (s *solrSearchDB) FindUser(queryString string) ([]*model.UserFound, error) {

	SolrQueryString := fmt.Sprintf("user_name:*%s* OR name:*%s*", queryString, queryString)
	query := solr.NewQuery()
	query.Q(SolrQueryString)
	search := s.Interface.Search(query)
	resp, err := search.Result(nil)
	if err != nil {
		return nil, err
	}
	var usersFound []*model.UserFound
	for _, doc := range resp.Results.Docs {
		userFound := s.docToUserFound(&doc)
		usersFound = append(usersFound, userFound)
	}
	return usersFound, nil

}
