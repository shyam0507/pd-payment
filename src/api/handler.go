package api

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shyam0507/pd-payment/src/types"
)

func (s *Server) updatePayment(c echo.Context) error {
	var orderId = c.Param("orderId")
	slog.Info("Updating payment for order", "OrderId", orderId)

	payment, err := s.storage.FindPayment(orderId)

	if err != nil || payment.Status == "COMPLETED" {
		slog.Error("Error while fetching the payment", "Err", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	payment.Status = "COMPLETED"

	if err := s.storage.UpdatePayment(payment.Id.Hex(), "COMPLETED"); err != nil {
		slog.Error("Error while updating the payment", "Err", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	s.producer.ProducePaymentReceivedEvent(payment.OrderId, types.PaymentReceivedEvent{
		Type:            "PaymentReceived",
		DataContentType: "application/json",
		Data:            types.PaymentInfo{OrderId: payment.OrderId, PaymentId: payment.Id.Hex(), Total: payment.Total},
	})

	return c.JSON(http.StatusOK, payment)
}
