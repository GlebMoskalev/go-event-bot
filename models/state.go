package models

type State string

const (
	StateStaffRegisterFullName    State = "state_staff_register_full_name"
	StateStaffRegisterPhoneNumber State = "state_staff_register_phone_number"
	StateStaffRegisterConfirm     State = "state_staff_register_confirm"
)
