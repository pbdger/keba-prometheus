package main

import (
	modbusclient "github.com/dpapathanasiou/go-modbus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

const metricsPort = "8080"

type environment struct {
	wallboxName string
	wallboxPort int
	debug       bool
}

type register struct {
	id          int
	description string
	value       int32
}

var (
	charingState = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "charging_state",
		Help:      "The register 1000 contains the state of the charging station.",
	})

	cableState = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "cable_state",
		Help:      "The register 1004 contains the state of the charging cable.",
	})

	errorCode = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "error_code",
		Help:      "The register 1006 contains the error code of the charging station.",
	})

	serialNumber = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "serial_number",
		Help:      "The register 1014 contains the serial number of the charging station.",
	})

	productTypeAndFeatures = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "product_type_and_features",
		Help:      "The register 1016 contains the type and the most important features of the charging station.",
	})

	firmwareVersion = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "firmware_version",
		Help: "The register 1018 contains the firmware version of the charging station." +
			" The represented value is in decimal and has to be converted to hex in order to read the firmware version.",
	})

	activePowerMW = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "active_power_mW",
		Help:      "The register 1020 contains the active power in milliwatts.",
	})

	totalEnergyCounterWh = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "total_energy_counter_Wh",
		Help:      "The register 1036 contains the active power in Wh.",
	})

	powerFactorPercent = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "power_factor_percent",
		Help:      "The register 1046 contains the current power factor (cos phi) in 0.1 %.",
	})

	maxChargingCurrentMAh = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "max_charging_current_mAh",
		Help:      "The register 1100 contains the maximum charging current of the charging station.",
	})
	maxSupportedCurrentMAh = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "max_supported_current_mAh",
		Help: "The register 1110 contains the maximum current value that can be supported by the hardware" +
			"of the charging station. This value represents the minimum of the DIP switch settings, cable coding" +
			" and temperature monitoring function.",
	})
	chargedEnergyWh = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "keba",
		Name:      "charged_energy_Wh",
		Help:      "The register 1502 contains the transferred energy of the current charging session.",
	})

	chargingCurrentPhase1MAh = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "keba",
		Name:        "charging_current_phase_mAh",
		Help:        "The register 1008,1010,1012 contains the measured current value on the phase in milliamperes.",
		ConstLabels: prometheus.Labels{"phase": "1"},
	})
	chargingCurrentPhase2MAh = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "keba",
		Name:        "charging_current_phase_mAh",
		Help:        "The register 1008,1010,1012 contains the measured current value on the phase in milliamperes.",
		ConstLabels: prometheus.Labels{"phase": "2"},
	})
	chargingCurrentPhase3MAh = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "keba",
		Name:        "charging_current_phase_mAh",
		Help:        "The register 1008,1010,1012 contains the measured current value on the phase in milliamperes.",
		ConstLabels: prometheus.Labels{"phase": "3"},
	})

	voltagePhase1V = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "keba",
		Name:        "voltage_phase_V",
		Help:        "The register 1040,1042,1044 contains the measured voltage value on the phase in volts.",
		ConstLabels: prometheus.Labels{"phase": "1"},
	})
	voltagePhase2V = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "keba",
		Name:        "voltage_phase_V",
		Help:        "The register 1040,1042,1044 contains the measured voltage value on the phase in volts.",
		ConstLabels: prometheus.Labels{"phase": "2"},
	})
	voltagePhase3V = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "keba",
		Name:        "voltage_phase_V",
		Help:        "The register 1040,1042,1044 contains the measured voltage value on the phase in volts.",
		ConstLabels: prometheus.Labels{"phase": "3"},
	})

	registers []register
	env       environment
)

func initRegisters() {

	log.Debug().Msg("Init registers")

	registers = append(registers, register{id: 1000, description: "charging_state", value: 0})
	registers = append(registers, register{id: 1004, description: "cable_state", value: 0})
	registers = append(registers, register{id: 1006, description: "error_code", value: 0})
	registers = append(registers, register{id: 1014, description: "serial_number", value: 0})
	registers = append(registers, register{id: 1016, description: "product_type_and_features", value: 0})
	registers = append(registers, register{id: 1018, description: "firmware_version", value: 0})
	registers = append(registers, register{id: 1020, description: "active_power_mW", value: 0})
	registers = append(registers, register{id: 1036, description: "total_energy_counter_Wh", value: 0})
	registers = append(registers, register{id: 1046, description: "power_factor_percent", value: 0})
	registers = append(registers, register{id: 1100, description: "max_charging_current_mAh", value: 0})
	registers = append(registers, register{id: 1110, description: "max_supported_current_mAh", value: 0})
	//registers = append(registers, register{id: 1502 ,description: "rfid_card", value: 0})
	registers = append(registers, register{id: 1502, description: "charged_energy_Wh", value: 0})

	registers = append(registers, register{id: 1008, description: "charging_current_phase_1_mAh", value: 0})
	registers = append(registers, register{id: 1010, description: "charging_current_phase_2_mAh", value: 0})
	registers = append(registers, register{id: 1012, description: "charging_current_phase_3_mAh", value: 0})
	registers = append(registers, register{id: 1040, description: "voltage_phase_1_V", value: 0})
	registers = append(registers, register{id: 1042, description: "voltage_phase_2_V", value: 0})
	registers = append(registers, register{id: 1044, description: "voltage_phase_3_V", value: 0})

}

func main() {

	initApp()

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		for {
			updateRegisterData()
			for _, register := range registers {
				log.Debug().Str(string(register.id), string(register.value))
				switch register.id {
				case 1000:
					charingState.Set(float64(register.value))
				case 1004:
					cableState.Set(float64(register.value))
				case 1006:
					errorCode.Set(float64(register.value))
				case 1014:
					serialNumber.Set(float64(register.value))
				case 1016:
					productTypeAndFeatures.Set(float64(register.value))
				case 1018:
					firmwareVersion.Set(float64(register.value))
				case 1020:
					activePowerMW.Set(float64(register.value))
				case 1036:
					totalEnergyCounterWh.Set(float64(register.value) / 10)
				case 1046:
					powerFactorPercent.Set(float64(register.value))
				case 1100:
					maxChargingCurrentMAh.Set(float64(register.value))
				case 1110:
					maxSupportedCurrentMAh.Set(float64(register.value))
				case 1502:
					chargedEnergyWh.Set(float64(register.value) / 10)
				case 1008:
					chargingCurrentPhase1MAh.Set(float64(register.value))
				case 1010:
					chargingCurrentPhase2MAh.Set(float64(register.value))
				case 1012:
					chargingCurrentPhase3MAh.Set(float64(register.value))
				case 1040:
					voltagePhase1V.Set(float64(register.value))
				case 1042:
					voltagePhase2V.Set(float64(register.value))
				case 1044:
					voltagePhase3V.Set(float64(register.value))
				}
			}
			time.Sleep(60 * time.Second)
		}
	}()

	log.Info().Msg("Find metrics on http://<your IP>:" + metricsPort + "/metrics (Not for usage in a container image!)")
	log.Fatal().Err(http.ListenAndServe(":"+metricsPort, nil))

}

func initApp() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	initEnvironmentVariables()

	if env.debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug level activated.")
	}

	initRegisters()
	initPrometheusRegisters()

}

func initPrometheusRegisters() {

	log.Debug().Msg("Init prometheus registers")
	prometheus.MustRegister(charingState)
	prometheus.MustRegister(cableState)
	prometheus.MustRegister(errorCode)
	prometheus.MustRegister(serialNumber)
	prometheus.MustRegister(productTypeAndFeatures)
	prometheus.MustRegister(firmwareVersion)
	prometheus.MustRegister(activePowerMW)
	prometheus.MustRegister(totalEnergyCounterWh)
	prometheus.MustRegister(powerFactorPercent)
	prometheus.MustRegister(maxChargingCurrentMAh)
	prometheus.MustRegister(maxSupportedCurrentMAh)
	prometheus.MustRegister(chargedEnergyWh)
	prometheus.MustRegister(chargingCurrentPhase1MAh)
	prometheus.MustRegister(chargingCurrentPhase2MAh)
	prometheus.MustRegister(chargingCurrentPhase3MAh)
	prometheus.MustRegister(voltagePhase1V)
	prometheus.MustRegister(voltagePhase2V)
	prometheus.MustRegister(voltagePhase3V)

}

func initEnvironmentVariables() {

	log.Info().Msg("Usage:")
	log.Info().Str("wallboxName", "This is an IP or a servername.").Msg("Mandatory environment parameter.")
	log.Info().Str("wallboxPort", strconv.Itoa(modbusclient.MODBUS_PORT)).Msg("Optional: The port TCP/modbus listens.")
	log.Info().Str("debug", "false").Msg("Optional: Use debug mode for logging (true | false ). ")

	env.wallboxName = getEnv("wallboxName", "")

	if len(env.wallboxName) == 0 {
		log.Fatal().Msg("The environment variable wallboxName is unset. Please fix this.")
	}

	portString := getEnv("wallboxPort", strconv.Itoa(modbusclient.MODBUS_PORT))
	port, err := strconv.Atoi(portString)
	env.wallboxPort = port

	if err != nil {
		log.Fatal().Err(err).Str("wallboxPort", portString)
	}

	debug := getEnv("debug", "false")
	env.debug, err = strconv.ParseBool(debug)

	if err != nil {
		log.Fatal().Err(err).Str("debug", debug)
	}

	log.Info().Str("wallboxName", env.wallboxName).Msg("This is the configured wallboxName.")
	log.Info().Str("wallboxPort", strconv.Itoa(env.wallboxPort)).Msg("This is the configured port TCP/modbus listens.")
	log.Info().Str("debug", strconv.FormatBool(env.debug)).Msg("Log debug mode.")

}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func updateRegisterData() {

	log.Debug().Msg("updateRegisterData")

	conn, err := modbusclient.ConnectTCP(env.wallboxName, env.wallboxPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Connection error.")
	}

	for i, register := range registers {
		log.Debug().Str("id", string(registers[i].id)).Str("value", string(registers[i].value))
		registers[i] = readRegister(conn, register)
		log.Debug().Str("id", string(registers[i].id)).Str("value", string(registers[i].value))
		time.Sleep(1 * time.Second)
	}

	modbusclient.DisconnectTCP(conn)

}

func readRegister(conn net.Conn, register register) register {

	// attempt to read one (0x01) holding registers starting at address 200
	readData := make([]byte, 3)
	readData[0] = byte(register.id >> 8)   // (High Byte)
	readData[1] = byte(register.id & 0xff) // (Low Byte)
	readData[2] = 0x01

	trace := zerolog.GlobalLevel() == zerolog.DebugLevel

	// make this read request transaction id 1, with a 300 millisecond tcp timeout
	readResult, readErr := modbusclient.TCPRead(conn, 300, 1, modbusclient.FUNCTION_READ_HOLDING_REGISTERS, false, 0x00, readData, trace)
	if readErr != nil {
		log.Fatal().Err(readErr)
	}

	var value int32
	var offset int

	from := len(readResult) - 4
	to := len(readResult)

	offset = 0
	value = 0

	for i := from; i < to; i++ {
		offset++
		switch offset {
		case 1:
			value = value + int32(readResult[i])*256*256*256
		case 2:
			value = value + int32(readResult[i])*256*256
		case 3:
			value = value + int32(readResult[i])*256
		case 4:
			value = value + int32(readResult[i])
		}
	}

	register.value = value

	return register
}
