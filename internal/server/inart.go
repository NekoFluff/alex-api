package server

import (
	"alex-api/data"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func (s *Server) InArt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		skip := (page - 1) * 50
		var limit int64 = 50

		log := logrus.New()
		db := data.New(log)
		twitterMediaList, err := db.GetTwitterMedia(&skip, &limit)

		if err != nil {
			l.Errorf("%v", err)
			w.Write([]byte(err.Error()))
			return
		}

		jsonStr, _ := json.MarshalIndent(twitterMediaList, "", "\t")
		l.Info("INART")
		l.Info(string(jsonStr))
		_, _ = w.Write(jsonStr)
	}
}
