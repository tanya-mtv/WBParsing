package api

import "fmt"

func prepareAPIBody(limit int64) string {
	return fmt.Sprintf(`{
        "sort": {
            "cursor": {
                "limit": %d
            },
            "filter": {
                "withPhoto": -1
            }
        }
      }`, limit)

}
