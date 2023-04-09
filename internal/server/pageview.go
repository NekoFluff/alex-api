package server

import (
	"alex-api/internal/data"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

func (s *Server) PageViewed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method":  r.Method,
			"path":    r.URL.EscapedPath(),
			"referer": r.Referer(),
		})

		url, err := url.Parse(r.Referer())
		if err != nil {
			l.Error("Invalid Referer URL")

			_, _ = w.Write([]byte("Invalid Referer URL"))
			return
		}

		pageView, err := s.db.GetPageView(url.Host, url.Path)
		if err != nil {
			pageView = data.PageView{
				Domain: url.Host,
				Path:   url.Path,
			}

			pageView.Increment()
			err = s.db.CreatePageView(pageView)
		} else {
			pageView.Increment()
			err = s.db.UpdatePageView(pageView)
		}

		if err != nil {
			l.Errorf("%v", err)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		jsonStr, _ := json.MarshalIndent(pageView, "", "\t")
		l.Info("Page View")
		_, _ = w.Write(jsonStr)
	}
}
