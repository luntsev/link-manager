package stat

import (
	"link-manager/configs"
	"link-manager/pkg/middleware"
	"link-manager/pkg/response"
	"log"
	"net/http"
	"time"
)

const (
	GroupByMonth = "month"
	GroupByDay   = "day"
)

type StatHandlerDeps struct {
	*StatRepository
	*configs.Config
}

type StatHandler struct {
	*StatRepository
}

func NewStatHendler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}

	router.Handle("GET /stat", middleware.IsAuthed(handler.getStat(), deps.Config))
}

func (handler *StatHandler) getStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			log.Println("wrong email in JWT token")
			return
		}

		from := r.URL.Query().Get("from")
		dateFrom, err := time.Parse("2006-01-02", from)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //400
			return
		}

		to := r.URL.Query().Get("to")
		dateTo, err := time.Parse("2006-01-02", to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //400
			return
		}

		period := r.URL.Query().Get("by")
		if period != GroupByDay && period != GroupByMonth {
			http.Error(w, "Invalid \"by\" param", http.StatusBadRequest) //400
			return
		}

		stats := handler.StatRepository.GetStats(period, dateFrom, dateTo)

		response.Json(w, stats, http.StatusOK)
	}
}
