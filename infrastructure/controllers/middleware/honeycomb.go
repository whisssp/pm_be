package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	honey "github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	baggagetrace "go.opentelemetry.io/contrib/processors/baggagecopy"
	"log"
	"os"
	"pm/infrastructure/config"
)

func HoneycombHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		envErr := godotenv.Load()
		if envErr != nil {
		}

		serviceName := "WHISPER"
		apiKey := config.GetEnv(fmt.Sprintf("honeycomb.%s.api_key", os.Getenv("ENV_NAME")), "local")
		// Enable multi-span attributes
		bsp := baggagetrace.NewSpanProcessor(baggagetrace.AllowAllMembers)

		// Use the Honeycomb distro to set up the OpenTelemetry SDK
		otelShutdown, err := otelconfig.ConfigureOpenTelemetry(
			otelconfig.WithSpanProcessor(bsp),
			otelconfig.WithServiceName(serviceName),
			honey.WithApiKey(apiKey),
		)
		if err != nil {
			log.Println("error setting up OTel SDK - %s", err)
		} else {
			//log.Println("Connected to honeycomb")
		}
		defer otelShutdown()
		// This is where you can put any custom logic you want to apply to all requests.
		// In this case, we're wrapping the request with OpenTelemetry instrumentation.
		otelgin.Middleware("gin-server")(c)
		c.Next()
	}
}