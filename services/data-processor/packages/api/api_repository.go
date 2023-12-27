package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/LukasJatmiko/simple-microservices/data-processor/driver"
	"github.com/LukasJatmiko/simple-microservices/data-processor/packages/sensor"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ApiRepository struct {
	DBDriver driver.Driver
}

func (ar *ApiRepository) GetSensorData(f *APIFilter) func(c echo.Context) error {
	return (func(c echo.Context) error {

		query := "select * from sensor_data where id1 is not null"
		var params []any

		//combo id
		f.ID1 = c.QueryParam("id1")
		f.ID2 = c.QueryParam("id2")
		if f.ID1 != "" && f.ID2 != "" {
			params = append(params, f.ID1, f.ID2)
			query += " and (id1 = ? and id2 = ?) "
		}

		//datetime range
		if c.QueryParam("starttime") != "" || c.QueryParam("endtime") != "" {
			if start, e := time.Parse("2006-01-02 15:04:05", c.QueryParam("starttime")); e != nil {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusBadRequest, Message: "bad request", Data: echo.Map{}})
			} else {
				f.start = start
			}
			if end, e := time.Parse("2006-01-02 15:04:05", c.QueryParam("endtime")); e != nil {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusBadRequest, Message: "bad request", Data: echo.Map{}})
			} else {
				f.end = end
			}
			if !f.end.After(f.start) {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusBadRequest, Message: "bad request", Data: echo.Map{}})
			} else {
				params = append(params, f.start.Format("2006-01-02 15:04:05"), f.end.Format("2006-01-02 15:04:05"))
				query += " and timestamp between ? and ? "
			}
		}

		//order data
		f.Order = strings.ToUpper(c.QueryParam("order"))
		if f.Order == "DESC" {
			query += " order by timestamp desc"
		}

		//pagination
		if c.QueryParam("page") != "" && c.QueryParam("rowsperpage") != "" {
			if page, e := strconv.ParseInt(c.QueryParam("page"), 10, 64); e != nil {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusInternalServerError, Message: "internal error", Data: echo.Map{}})
			} else {
				f.Page = int(page)
			}
			if rpp, e := strconv.ParseInt(c.QueryParam("rowsperpage"), 10, 64); e != nil {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusInternalServerError, Message: "internal error", Data: echo.Map{}})
			} else {
				f.RowPerPage = int(rpp)
			}
			if f.Page > 0 && f.RowPerPage > 0 {
				params = append(params, ((f.Page - 1) * f.RowPerPage), (f.Page * f.RowPerPage))
				query += " limit ?, ?"
			}
		}

		result := ar.DBDriver.GetWrapperInstance().(*gorm.DB).Raw(query, params...)
		if result.Error != nil {
			return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusInternalServerError, Message: "internal error", Data: echo.Map{}})
		} else {
			var sensorData []sensor.SensorData
			if e := result.Scan(&sensorData).Error; e != nil {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusInternalServerError, Message: "internal error", Data: echo.Map{}})
			} else {
				if len(sensorData) < 1 {
					return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusOK, Message: "ok", Data: make([]int, 0)})
				} else {
					return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusOK, Message: "ok", Data: sensorData})
				}
			}
		}
	})
}

func (ar *ApiRepository) DeleteSensorData(f *APIFilter) func(c echo.Context) error {
	return (func(c echo.Context) error {

		query := "delete from sensor_data where"
		var params []any

		//combo id
		f.ID1 = c.QueryParam("id1")
		f.ID2 = c.QueryParam("id2")
		if f.ID1 != "" && f.ID2 != "" {
			params = append(params, f.ID1, f.ID2)
			query += " (id1 = ? and id2 = ?) and"
		}

		//datetime range
		if c.QueryParam("starttime") != "" || c.QueryParam("endtime") != "" {
			if start, e := time.Parse("2006-01-02 15:04:05", c.QueryParam("starttime")); e != nil {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusBadRequest, Message: "bad request", Data: echo.Map{}})
			} else {
				f.start = start
			}
			if end, e := time.Parse("2006-01-02 15:04:05", c.QueryParam("endtime")); e != nil {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusBadRequest, Message: "bad request", Data: echo.Map{}})
			} else {
				f.end = end
			}
			if !f.end.After(f.start) {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusBadRequest, Message: "bad request", Data: echo.Map{}})
			} else {
				params = append(params, f.start.Format("2006-01-02 15:04:05"), f.end.Format("2006-01-02 15:04:05"))
				query += " timestamp between ? and ? and"
			}
		}

		//remove trailing and
		query = query[:len(query)-3] + ";"

		result := ar.DBDriver.GetWrapperInstance().(*gorm.DB).Exec(query, params...)
		if result.Error != nil {
			return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusInternalServerError, Message: "internal error", Data: echo.Map{}})
		} else {
			return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusOK, Message: fmt.Sprintf("%v rows successfully deleted", result.RowsAffected), Data: echo.Map{}})
		}
	})
}

func (ar *ApiRepository) UpdateSensorData(f *APIFilter) func(c echo.Context) error {
	return (func(c echo.Context) error {

		query := "update sensor_data set value = ? where"
		var params []any

		if nv, e := strconv.ParseFloat(c.QueryParam("newvalue"), 64); e != nil {
			return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusBadRequest, Message: "bad request", Data: echo.Map{}})
		} else {
			params = append(params, nv)
		}

		//combo id
		f.ID1 = c.QueryParam("id1")
		f.ID2 = c.QueryParam("id2")
		if f.ID1 != "" && f.ID2 != "" {
			params = append(params, f.ID1, f.ID2)
			query += " (id1 = ? and id2 = ?) and"
		}

		//datetime range
		if c.QueryParam("starttime") != "" || c.QueryParam("endtime") != "" {
			if start, e := time.Parse("2006-01-02 15:04:05", c.QueryParam("starttime")); e != nil {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusBadRequest, Message: "bad request", Data: echo.Map{}})
			} else {
				f.start = start
			}
			if end, e := time.Parse("2006-01-02 15:04:05", c.QueryParam("endtime")); e != nil {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusBadRequest, Message: "bad request", Data: echo.Map{}})
			} else {
				f.end = end
			}
			if !f.end.After(f.start) {
				return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusBadRequest, Message: "bad request", Data: echo.Map{}})
			} else {
				params = append(params, f.start.Format("2006-01-02 15:04:05"), f.end.Format("2006-01-02 15:04:05"))
				query += " timestamp between ? and ? and"
			}
		}

		//remove trailing and
		query = query[:len(query)-3] + ";"

		result := ar.DBDriver.GetWrapperInstance().(*gorm.DB).Exec(query, params...)
		if result.Error != nil {
			return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusInternalServerError, Message: "internal error", Data: echo.Map{}})
		} else {
			return c.JSON(http.StatusOK, APIResponsePayload{Status: http.StatusOK, Message: fmt.Sprintf("%v rows successfully updated", result.RowsAffected), Data: echo.Map{}})
		}
	})
}
