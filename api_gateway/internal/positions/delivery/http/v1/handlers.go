package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go-futures-api/internal/models"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type PositionsHandlers struct {
	group  *gin.RouterGroup
	logger *zap.SugaredLogger
}

func NewPositionsHandlers(group *gin.RouterGroup, logger *zap.SugaredLogger) *PositionsHandlers {
	return &PositionsHandlers{
		group:  group,
		logger: logger,
	}
}

// Register GetPositionDetail
// @Tags Positions
// @Summary Get position detail by id
// @Description Get single position by id
// @Accept json
// @Produce json
// @Param position_id query integer true "position id"
// @Success 200 {object} models.Position
// @Router /positions/{position_id} [get]
func (p *PositionsHandlers) GetPositionDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "PositionsHandlers.GetPositionDetail")
		fmt.Println(ctx)
		defer span.Finish()

		id, err := strconv.Atoi(c.Query("position_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, &models.Position{PositionId: int64(id)})
	}
}

func (p *PositionsHandlers) MapRoutes() {
	p.group.GET("/:position_id", p.GetPositionDetail())
}
