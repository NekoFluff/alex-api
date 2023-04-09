package server

import (
	"alex-api/internal/data"
	"alex-api/internal/utils"
	"encoding/json"
	"io"
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

		path := url.Path

		if r.Body != http.NoBody {
			var body data.PageView
			err = utils.DecodeValidate(r.Body, s.validator, &body)
			defer r.Body.Close()
			if err != nil && err != io.EOF {
				l.WithError(err).Error("failed to decode request")
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			l = l.WithField("body", body)

			if body.Path != "" {
				path = body.Path
			}
		}

		pageView, err := s.db.GetPageView(url.Host, path)
		if err != nil {
			pageView = data.PageView{
				Domain: url.Host,
				Path:   path,
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
