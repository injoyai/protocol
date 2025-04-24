package j60

const (
	// Commands
	DISABLE_MOTOR         uint8 = 1
	ENABLE_MOTOR          uint8 = 2
	CALIBRATE_START       uint8 = 3
	CONTROL_MOTOR         uint8 = 4
	RESET_MOTOR           uint8 = 5
	SET_HOME              uint8 = 6
	SET_GEAR              uint8 = 7
	SET_ID                uint8 = 8
	SET_CAN_TIMEOUT       uint8 = 9
	SET_BANDWIDTH         uint8 = 10
	SET_LIMIT_CURRENT     uint8 = 11
	SET_UNDER_VOLTAGE     uint8 = 12
	SET_OVER_VOLTAGE      uint8 = 13
	SET_MOTOR_TEMPERATURE uint8 = 14
	SET_DRIVE_TEMPERATURE uint8 = 15
	SAVE_CONFIG           uint8 = 16
	ERROR_RESET           uint8 = 17
	WRITE_APP_BACK_START  uint8 = 18
	WRITE_APP_BACK        uint8 = 19
	CHECK_APP_BACK        uint8 = 20
	DFU_START             uint8 = 21
	GET_FW_VERSION        uint8 = 22
	GET_STATUS_WORD       uint8 = 23
	GET_CONFIG            uint8 = 24
	CALIB_REPORT          uint8 = 31
)
